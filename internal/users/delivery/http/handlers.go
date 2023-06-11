package http

import (
	"UsersService/config"
	"UsersService/internal/users"
	"UsersService/pkg/auth"
	"UsersService/pkg/cdnService"
	"UsersService/pkg/logger"
	"UsersService/pkg/loggerService"
	"UsersService/pkg/rating"
	"UsersService/pkg/sso"
	"UsersService/pkg/utils"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UsersHandler struct {
	logger   *logger.ApiLogger
	cfg      *config.Config
	userUC   users.UseCase
	ratingSR rating.ServiceRating
	authSR   auth.ServiceAuth
	ssoSR    sso.ServiceSso
	cdn      cdnService.UseCase
}

func NewUsersHandlers(cfg *config.Config, logger *logger.ApiLogger, userUC users.UseCase, ratingSR rating.ServiceRating,
	authSR auth.ServiceAuth, ssoSR sso.ServiceSso, cdn cdnService.UseCase) *UsersHandler {
	return &UsersHandler{
		cfg:      cfg,
		logger:   logger,
		userUC:   userUC,
		ratingSR: ratingSR,
		authSR:   authSR,
		ssoSR:    ssoSR,
		cdn:      cdn,
	}
}

func (u *UsersHandler) ClientSignInTG() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.ClientSignInTGParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}

		data := u.userUC.ClientSignInTG(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент %s не смог зайти через tg", params.TgUserName), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ClientSignInTG -> data = {%+v}", params), 4)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) ClientSignIn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.SignInParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		params.UA = ctx.Get("ua")
		if params.UA == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.ClientSignIn(&users.SignInSAParams{Email: params.Email, Nickname: params.Nickname,
			Password: params.Password, UA: params.UA, Ip: params.Ip})
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			params.Password = "********"
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент %s не смог зайти {%+v}", params.Nickname, params), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ClientSignIn -> data = {%+v}", params), 4)
		}
		if data.Error.InternalCode == 0 {
			go u.userUC.UpdateLastLogin(&users.UpdateLastLoginParams{Email: params.Email, Ip: params.Ip})
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) ClientSignUp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.SignUpParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		params.UA = ctx.Get("ua")
		data := u.userUC.ClientSignUp(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			params.Password = "*********"
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент не смог зарегистрироваться, используя следующие параметры: {%+v}, Ошибка: %s", &params, &data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ClientSignIn -> data = {%+v}, error = {%+v}", &params, &data.Error), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) Logout() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userUUID := u.getUserUUID(ctx)
		UA := ctx.Get("ua")
		data := u.ssoSR.Logout(&sso.LogoutParams{ClientUUID: userUUID, UA: UA})
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент id: %s не смог разлогиниться", userUUID), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> Logout -> data = {%+v}", &data.Error), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) StopSession() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userUUID := u.getUserUUID(ctx)
		var params users.UaParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.ssoSR.Logout(&sso.LogoutParams{
			ClientUUID: userUUID,
			UA:         params.UA,
		})
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент id: %s не смог разлогиниться, причина: %s", userUUID, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> Logout -> data = {%+v}", &data.Error), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) CreateClientDND() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.CreateClientParams
		data := u.userUC.CreateClientDND(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Мы не смогли создать dnd client"), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> CreateClient -> data = {%+v} error = {%+v}", params, &data.Error), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) ChangePassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.ChangePasswordParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		params.Access = ctx.Get("Access")
		fmt.Println(params.Totp)
		user := u.userUC.GetClientUUIDByAccessToken(&users.GetUserIDByAccessTokenParams{Access: params.Access})
		if user.Error.InternalCode != 0 {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		data := u.userUC.ClientChangePassword(&users.ChangeClientPasswordParams{
			ClientUuid:       user.Data.ClientUUID,
			OldPassword:      params.OldPassword,
			NewPassword:      params.NewPassword,
			NewPasswordAgain: params.NewPasswordAgain,
			Totp:             params.Totp,
		})
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %s не смог сменить пароль err: %s", user.Data.ClientUUID, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"Users -> ChangePassword -> data = {%+v}, err = {%+v}", &params, &data.Error), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) IsUserBlocked() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.IsUserBlockedParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.userUC.IsUserBlocked(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}

}

func (u *UsersHandler) GetClientRating() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserID := u.getUserID(ctx)
		data := u.userUC.GetClientRating(&users.GetClientRatingParams{ClientID: UserID})
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}

}

func (u *UsersHandler) GetUserIDByAccessToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserIDByAccessTokenParams
		params.Access = ctx.Get("Access", "")
		if params.Access == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.GetClientUUIDByAccessToken(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserIDByAccessToken -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) GetUserByAccessToken() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserByAccessTokenParams
		params.Token = ctx.Get("Access", "")
		if params.Token == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.GetUserByAccessToken(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}

		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserByAccessToken -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}

}

func (u *UsersHandler) GetUserSessions() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserIDByAccessTokenParams
		params.Access = ctx.Get("Access", "")
		if params.Access == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		user := u.userUC.GetClientUUIDByAccessToken(&params)
		if user.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", user.Error.Message, user.Error.InternalCode)
		}

		data := u.ssoSR.GetUserSessions(&sso.GetUserSessions{ClientUuid: user.Data.ClientUUID})
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"🙀Клиент %s не смог получить свои сессии err: %s", user.Data.ClientUUID, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"Users -> GetUserSessions -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}

}

func (u *UsersHandler) GetUserByAccessTokenWithID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserByAccessTokenParams
		params.Token = ctx.Get("Access", "")
		if params.Token == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.GetUserByAccessTokenWithID(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}

		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}

}

func (u *UsersHandler) Recovery() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.RecoveryParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.userUC.Recovery(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент не смог инициализировать смену пароля err: %s (забыл пароль) email: %s", data.Error.Message, params.Email), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> Recovery -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)

	}
}

func (u *UsersHandler) RecoveryPassword() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.RecoveryConfirmParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		params.Hash = ctx.Params("hash", "")
		if params.Hash == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.RecoveryConfirm(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент не смог восстановить аккаунт смена пароля err: %s hash: %s", data.Error.Message, params.Hash), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> RecoveryPassword -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)

	}
}

func (u *UsersHandler) ConfirmEmail() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.ConfirmEmailParams
		params.Access = ctx.Get("Access", "")
		if params.Access == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.ConfirmEmail(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент не смог инициализировать подтверждение почты err: %s", data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ConfirmEmail -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)

	}
}

func (u *UsersHandler) AddTotp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserIDByAccessTokenParams
		params.Access = ctx.Get("Access", "")
		if params.Access == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		_user := u.userUC.GetClientUUIDByAccessToken(&params)
		if _user.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", _user.Error.Message, _user.Error.InternalCode)
		}

		data := u.ssoSR.AddTotp(&sso.AddTotpParams{ClientUuid: _user.Data.ClientUUID})
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"🙀Клиент %s не смог добавить 2fa err: %s", _user.Data.ClientUUID, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> AddTotp -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)

	}
}

func (u *UsersHandler) KycInit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params sso.KycInitParams
		params.ClientUuid = u.getUserUUID(ctx)
		clientId := u.getUserID(ctx)
		if params.ClientUuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.ssoSR.ConfirmKycInit(&params)
		if data.Error.InternalCode != 0 {
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> KycInit -> params = {%+v}, data = {%+v}", &params, &data), 8)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %s не смог инициализировать kyc err: %s", params.ClientUuid, data.Error.Message), 8)
			return ctx.Status(int(data.Error.StandartCode)).JSON(data)
		}

		data2 := u.userUC.ConfirmKyc(&users.HashKycConfirm{Hash: data.Data.Hash})
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"🙀Клиент %s инициализировал kyc", params.ClientUuid), 8)
		if data2.Error.InternalCode != 0 && data2.Error.InternalCode != 6054 {
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"Users -> KycInit -> params = {%+v}, data = {%+v}", &params, &data), 8)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %s не смог инициализировать kyc err: %s", params.ClientUuid, data.Error.Message), 8)
			return ctx.Status(int(data2.Error.StandartCode)).JSON(data2)
		}
		u.userUC.AddAuthKyc(&users.AuthKycParams{AuthToken: data2.Data.AuthToken, ClientID: clientId})
		return ctx.JSON(data2)

	}
}

func (u *UsersHandler) KycConfirm() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params sso.KycParams
		var body users.KYCResponse
		if err := utils.ReadRequest(ctx, &body); err != nil {
			badBody := string(ctx.Body())
			go loggerService.GetInstance().DevLog(fmt.Sprintf("!🙀! CallBack от idefi BADREQUEST:\n %+v", badBody), 8)
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		body.FileUrls.BACK = "****"
		body.FileUrls.FRONT = "****"
		body.FileUrls.FACE = "***"
		go loggerService.GetInstance().DevLog(fmt.Sprintf("! 🙀! CallBack от idefi: %+v", body), 8)
		if body.Status.Overall != "APPROVED" {
			go loggerService.GetInstance().DevLog(fmt.Sprintf("!🙀! CallBack от idefi not APPROVED: %+v", body), 8)
			return ctx.SendStatus(200)
		}
		params.Hash = body.ClientId
		data := u.ssoSR.ConfirmKyc(&params)
		if data.Error.InternalCode == 0 {
			go loggerService.GetInstance().DevLog(fmt.Sprintf("!🙀! CallBack от idefi BADREQUEST: %+v", &data), 8)
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)

	}
}

//func (u *UsersHandler) KycConfirmErr() fiber.Handler {
//	return func(ctx *fiber.Ctx) error {
//		go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Пришел callback ошибки от idefi: %+v", ctx), 1)
//		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> KycInit -> params = {%+v}, data = {%+v}", &params, &data), 4)
//		return ctx.SendStatus(200)
//	}
//}

func (u *UsersHandler) VerifyTotp() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserIDByAccessTokenParams
		params.Access = ctx.Get("Access", "")
		if params.Access == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		_user := u.userUC.GetClientUUIDByAccessToken(&params)
		if _user.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", _user.Error.Message, _user.Error.InternalCode)
		}

		var params2 sso.VerifyTotpParams
		if err := utils.ReadRequest(ctx, &params2); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		data := u.ssoSR.VerifyTotp(&sso.VerifyTotpParams{
			ClientUuid: _user.Data.ClientUUID,
			Token:      params2.Token,
		})
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"🙀Клиент %s не смог верифицировать 2fa err: %s", _user.Data.ClientUUID, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"Users -> VerifyTotp -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)

	}

}

func (u *UsersHandler) VerifyTotpInit() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserIDByAccessTokenParams
		params.Access = ctx.Get("Access", "")
		if params.Access == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		_user := u.userUC.GetClientUUIDByAccessToken(&params)
		if _user.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", _user.Error.Message, _user.Error.InternalCode)
		}

		var params2 sso.VerifyTotpParams
		if err := utils.ReadRequest(ctx, &params2); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		data := u.ssoSR.VerifyTotpInit(&sso.VerifyTotpParams{
			ClientUuid: _user.Data.ClientUUID,
			Token:      params2.Token,
		})
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"🙀Клиент %s не смог инициализировать добавление 2fa err: %s", _user.Data.ClientUUID, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"Users -> VerifyTotpInit -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)

	}

}

func (u *UsersHandler) VerifyTotpWithID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {

		var params sso.VerifyTotpParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}

		data := u.ssoSR.VerifyTotp(&sso.VerifyTotpParams{
			ClientUuid: params.ClientUuid,
			Token:      params.Token,
		})
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> VerifyTotpWithID -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}

}

func (u *UsersHandler) GetUserByNickname() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserByNickNameParams
		params.Nickname = ctx.Params("nickname")
		data := u.userUC.GetUserByNickName(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserByNickname -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) GetUserByNicknameWithID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserByNicknameWithID
		params.Nickname = ctx.Params("nickname")
		data := u.userUC.GetUserByNicknameWithID(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserByNicknameWithID -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) GetUserReviewsByNickname() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserByNickNameParams
		queryString := string(ctx.Request().URI().QueryString())
		query := utils.ParseQuery(queryString)
		params.Nickname = ctx.Params("nickname")
		params.Type = query["type"]
		lim, err := strconv.ParseInt(query["limit"], 10, 64)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		params.Limit = lim
		pag, err := strconv.ParseInt(query["page"], 10, 64)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		params.Page = pag
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.userUC.GetUserReviewsByNickname(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент не смог посмотреть отзывы  %s err: %s", params.Nickname, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserReviewsByNickname -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) GetClientStatistics() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserUUID := u.getUserID(ctx)
		data := u.ratingSR.GetClientSRRating(&rating.GetClientSRRatingParams{ClientID: UserUUID})
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetClientStatistics -> err = {%+v}", data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) GetRegistration() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var preParams users.GetRegistrationStringParams
		var params users.GetRegistrationParams
		utils.ReadRequest(ctx, &preParams)

		params.FromDate = utils.ParseStringToTime(preParams.FromDate)
		params.ToDate = utils.ParseStringToTime(preParams.ToDate)
		params.Dnd = preParams.Dnd

		fmt.Println(params)
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.userUC.GetRegistration(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetClientStatistics -> err = {%+v}", data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) AddComment() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserID := u.getUserID(ctx)
		var params rating.AddCommentParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		params.ReviewerID = UserID
		params.CreatedDate = utils.GetEuropeTime()
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.ratingSR.AddComment(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент %d не смог осативть отзыв err: %s", UserID, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> VerifyTotpInit -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) IsCommentExist() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params rating.IsCommentExistParams
		UserID := u.getUserID(ctx)
		params.InternalID = ctx.Params("internal_id", "")
		params.ClientID = UserID
		data := u.ratingSR.IsCommentExist(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> IsCommentExist -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) UpdateCommentClient() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserID := u.getUserID(ctx)
		var params rating.UpdateCommentClientParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		params.ClientID = UserID
		data := u.ratingSR.UpdateCommentClient(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент %d не смог изменить комментарий err: %s", UserID, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> UpdateCommentClient -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) UpdateCommentAdmin() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params *rating.UpdateCommentAdminParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.ratingSR.UpdateCommentAdmin(params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> UpdateCommentAdmin -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) UpdateUserBio() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserID := u.getUserID(ctx)
		var params users.UpdateUserBioParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		params.ClientID = UserID
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.userUC.UpdateUserBio(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"🙀Клиент %d не смог оставить отзыв err: %s", UserID, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"Users -> VerifyTotpInit -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) UpdateUserNickName() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.UpdateUserNickNameParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		params.ClientID = u.getUserID(ctx)
		params.Refresh = ctx.Get("Refresh")
		params.UA = ctx.Get("UA")
		data := u.userUC.UpdateUserNickName(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %s не смог изменить nickname err: %s", params.ClientID, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"Users -> UpdateUserNickName -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) UpdateUserAvatar() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.UpdateUserAvatarParams
		params.ClientID = u.getUserID(ctx)
		file, err := ctx.FormFile("image")
		if err != nil {
			loggerService.GetInstance().DevLog(fmt.Sprintf("UpdateUserAvatar -> {%s}", err.Error()), 4)
			return ctx.JSON(fiber.StatusBadRequest)
		}
		src, err := file.Open()
		if err != nil {
			return ctx.JSON(fiber.StatusBadRequest)
		}

		img, err := ioutil.ReadAll(src)
		if err != nil {
			loggerService.GetInstance().DevLog(fmt.Sprintf("Order -> SaveImg -> {%s}", err.Error()), 4)
			return ctx.JSON(fiber.StatusInternalServerError)
		}
		base64Encoding := base64.StdEncoding.EncodeToString(img)

		dataImage, err := u.cdn.SaveImage(&cdnService.SaveImageParams{ImageBase64: base64Encoding})
		if err != nil {
			loggerService.GetInstance().DevLog(fmt.Sprintf("Order -> SaveImg -> {%s}", err.Error()), 4)
			return ctx.JSON(fiber.StatusInternalServerError)

		}
		params.NewAvatar = dataImage.File
		data := u.userUC.UpdateUserAvatar(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %d не смог изменить аватар err: %s", params.ClientID, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"Users -> UpdateUserAvatar -> data = {%+v}, err = {%+v}", &params, &data), 4)
			return ctx.Status(int(data.Error.StandartCode)).JSON(data)
		}
		return ctx.JSON(data)
	}
}

func (u *UsersHandler) ChangeLanguage() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserIDByAccessTokenParams
		params.Access = ctx.Get("Access", "")
		if params.Access == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		data := u.userUC.GetUserIDByAccessToken(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserIDByAccessToken -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		var params2 users.ChangeDefaultUserLanguageParams
		if err := utils.ReadRequest(ctx, &params2); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		params2.ClientID = data.Data.ClientID
		data2 := u.userUC.ClientChangeDefaultLanguage(&params2)
		if data2.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data2.Error.Message, data2.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %d не смог изменить язык err: %s", data.Data.ClientID, data2.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"Users -> UpdateUserAvatar -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data2.Error.InternalCode == 0 {
			return ctx.JSON(data2)
		}
		return ctx.Status(int(data2.Error.StandartCode)).JSON(data2)
	}
}

func (u *UsersHandler) ChangeTg() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params sso.ChangeTgParams
		userUUID := u.getUserUUID(ctx)
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		if err := utils.ValidateStructSize(params); err != nil {
			return ctx.SendStatus(fiber.StatusRequestEntityTooLarge)
		}
		params.ClientUuid = userUUID
		data := u.ssoSR.ChangeTg(&params)
		if data.Error.InternalCode != 0 {
			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %s не смог изменить tg err: %s", userUUID, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"Users -> UpdateUserAvatar -> data = {%+v}, err = {%+v}", &params, &data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) getUserID(ctx *fiber.Ctx) int64 {
	var params users.GetUserIDByAccessTokenParams
	params.Access = ctx.Get("Access", "")
	if params.Access == "" {
		return 0
	}
	data := u.userUC.GetUserIDByAccessToken(&params)
	if data.Error.InternalCode != 0 {
		u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
		// TODO: tg
	}
	return data.Data.ClientID
}

func (u *UsersHandler) getUserUUID(ctx *fiber.Ctx) string {
	var params users.GetUserIDByAccessTokenParams
	params.Access = ctx.Get("Access", "")
	if params.Access == "" {
		return ""
	}
	data := u.userUC.GetClientUUIDByAccessToken(&params)
	if data.Error.InternalCode != 0 {
		u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
		// TODO: tg
	}
	return data.Data.ClientUUID
}

func (u *UsersHandler) GetUserByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userId, err := strconv.Atoi(ctx.Params("user_id"))
		if err != nil {
			// TODO: tg
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		data := u.userUC.GetUserByID(&users.GetUserByIDParams{ClientID: int64(userId)})
		if data.Error.InternalCode != 0 {

			u.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserByID -> data = {%+v}", data), 4)
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) ConfirmEmailByHash() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.ConfirmEmailByHashParams
		params.Hash = ctx.Params("hash")
		data := u.userUC.ConfirmEmailByHash(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент не смог подтвердить почту hash: %s err: %s", params.Hash, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ConfirmEmailByHash -> data = {%+v}, err = {%+v}", &params, &data), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) SendCodeToEmail(typeCode string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params sso.SendCodeToEmailParams
		params.ClientUuid = u.getUserUUID(ctx)
		userId := u.getUserID(ctx)
		user := u.userUC.GetUserByID(&users.GetUserByIDParams{ClientID: userId})
		params.Type = typeCode
		params.LanguageIso = user.Data.Language
		data := u.ssoSR.SendCodeToEmail(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) CheckCodeFromEmail(typeCode string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params sso.VerCodeParams
		params.Type = typeCode
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := u.ssoSR.CheckCodeFromEmail(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) RefreshAccess() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params sso.RefreshAccessParams
		params.Refresh = ctx.Get("Refresh", "")
		params.UA = ctx.Get("ua", "")
		params.ClientUuid = u.getUserUUID(ctx)
		if params.Refresh == "" || params.UA == "" || params.ClientUuid == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.ssoSR.RefreshAccess(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) ValidAccess() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserUUID := u.getUserUUID(ctx)
		var params sso.ValidateAccessTokenParams
		params.Access = ctx.Get("Access", "")
		params.UA = ctx.Get("ua", "")
		if params.Access == "" || params.UA == "" {
			data := &sso.ResponseValidateAccess{
				Data: &sso.ResponseSuccessModel{
					Success: true,
					Message: "success but token or ua is empty",
				},
			}
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %d не смог верифицировать access, c фронта пришел пустой токен params = {%+v}", UserUUID, params), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ValidAccess -> data = {%+v}, err = {%+v}", &params, &data), 4)
			return ctx.JSON(data)
		}
		if params.Access == "" || params.UA == "" || UserUUID == "" {
			go loggerService.GetInstance().DevLog(fmt.Sprintf(
				"🙀Клиент %s не смог верифицировать access, c фронта пришел пустой токен params = {%+v}", UserUUID, params), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ValidAccess -> data = {%+v}", params), 4)
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		params.ClientUuid = UserUUID
		data := u.ssoSR.ValidateAccess(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"🙀Клиент %s не смог верифицировать access, params = {%+v}, err = %s", UserUUID, &params, data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf(
			"Users -> ValidAccess -> data = {%+v}, params = {%+v}", data, params), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) ReadNotice() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserID := u.getUserID(ctx)
		var params users.ReadNoticeParams
		params.ClientID = UserID
		params.InternalID = ctx.Params("internal_id", "")
		if params.ClientID == 0 || params.InternalID == "" {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.ReadNotice(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> ReadNotice -> data = {%+v}, params = {%+v}", data, params), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) GetActiveNotice() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		UserID := u.getUserID(ctx)
		var params users.GetActiveNoticeParams
		params.ClientID = UserID
		params.TypeNotice = ctx.Query("type", "")
		if params.ClientID == 0 || params.TypeNotice == "" {
			fmt.Println(params.TypeNotice)
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}
		data := u.userUC.GetActiveNotice(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент %d не смог получить уведомения, params = {%+v}, err = %s", UserID, &params, &data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetActiveNotice -> data = {%+v}, params = {%+v}", data, params), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) AddNotice() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.ADDNotice
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := u.userUC.AddNotice(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент не смог получить уведомения, params = {%+v}, err = %s", &params, &data.Error.Message), 1)
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> AddNotice -> data = {%+v}, params = {%+v}", data, params), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (u *UsersHandler) GetUserNicknameByID() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.GetUserNicknameByIDParams
		params.ClientID = utils.StringToInt(ctx.Params("client_id"))
		data := u.userUC.GetUserNicknameByID(&params)
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> GetUserNicknameByID -> data = {%+v}, params = {%+v}", data, params), 4)
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}
