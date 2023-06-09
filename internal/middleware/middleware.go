package middleware

import (
	"UsersService/config"
	"UsersService/internal/admins"
	"UsersService/internal/cConstants"
	"UsersService/internal/cErrors"
	"UsersService/internal/iConnection"
	"UsersService/internal/users"
	"UsersService/pkg/auth"
	"UsersService/pkg/logger"
	"UsersService/pkg/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

type MDWManager struct {
	cfg           *config.Config
	usersUC       users.UseCase
	adminUC       admins.AdminCase
	authSR        auth.ServiceAuth
	logger        logger.ApiLogger
	iConnectionUC iConnection.UseCase
}

func NewMDWManager(cfg *config.Config, usersUC users.UseCase, adminUC admins.AdminCase, authSR auth.ServiceAuth, iConnectionUC iConnection.UseCase, logger logger.ApiLogger) *MDWManager {
	return &MDWManager{usersUC: usersUC, cfg: cfg, adminUC: adminUC, authSR: authSR, iConnectionUC: iConnectionUC, logger: logger}
}

func (mw *MDWManager) VerifySignatureMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var (
			publicKey string = c.Get("ApiPublic")
			signature string = c.Get("Signature")
			timestamp string = c.Get("Timestamp")
			method    string = c.Method()
			err       error
			message   string
		)

		if method == "POST" {
			message = string(c.Body())
			message = strings.ReplaceAll(message, "\n", "")
		} else {
			message = publicKey
		}
		// ------------------------------------------------------------------------------------------------

		if timestamp == "" {
			err = errors.New("timestamp is required")
			mw.logger.Error(err)
			return err
		}

		isValid, err := mw.iConnectionUC.Validate(&iConnection.ValidateParams{
			Signature: signature,
			Public:    publicKey,
			Message:   message,
			Timestamp: timestamp,
		})
		if err != nil || !*isValid {
			mw.logger.Error(err)
			return err
		}
		mw.logger.Info("OK")
		return c.Next()
	}
}

func (mw *MDWManager) UpdateUserLastActivity() fiber.Handler {
	return func(c *fiber.Ctx) error {

		access := c.Get("Access", "")
		if access == "" {
			return c.JSON(fiber.StatusUnauthorized)
		}
		user := mw.usersUC.GetClientUUIDByAccessToken(&users.GetUserIDByAccessTokenParams{Access: access})
		if user.Error.InternalCode != 0 {
			//fmt.Println(user.Data.ClientID, "!!!!!!!!!!!!!!!!!!!!!")
		}
		//res := mw.usersUC.UpdateUserLastActivity(&users.UpdateUserLastActivityParams{ClientID: user.Data.ClientID})
		//if res.Error.InternalCode != 0 {
		//	mw.logger.Info(res.Error.Message)
		//}
		return c.Next()
	}
}

func (mw *MDWManager) ValidateToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params auth.ValidateAccessTokenParams
		params.Access = c.Get("Access")
		if params.Access == "" {
			res := &users.ResponseSuccessModel{Success: false,
				Message: "not valid (empty) Access"}
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}
		params.UA = c.Get("ua")
		if params.UA == "" {
			return c.Next()
		}
		client := mw.usersUC.GetUserByAccessTokenWithID(&users.GetUserByAccessTokenParams{Token: params.Access})
		if client.Error.InternalCode != 0 {
			res := &users.ResponseSuccessModel{Success: false,
				Message: "not valid (no such user) Access"}
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}
		params.ClientID = client.Data.ClientID
		valid := mw.authSR.ValidateAccess(&params)
		if valid.Error.InternalCode != 0 {
			res := &users.ResponseSuccessModel{Success: false,
				Message: "not valid (error) Access"}
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}
		if !valid.Data.Success {
			res := &users.ResponseSuccessModel{Success: false,
				Message: "not valid Access"}
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}
		return c.Next()
	}
}

func (mw *MDWManager) ValidateVerifications() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var params auth.ValidateAccessTokenParams
		params.Access = c.Get("Access")
		if params.Access == "" {
			res := &users.ResponseSuccessModel{Success: false,
				Message: "not valid (empty) Access"}
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}
		client := mw.usersUC.GetClientUUIDByAccessToken(&users.GetUserIDByAccessTokenParams{Access: params.Access})
		if client.Error.InternalCode != 0 {
			res := &users.ResponseSuccessModel{Success: false,
				Message: "not valid (no such user) Access"}
			return c.Status(fiber.StatusUnauthorized).JSON(res)
		}
		check, _ := strconv.Atoi(client.Data.ClientUUID) // строки не было
		params.ClientID = int64(check)                   // client.Data.ClientUUID было 07.06 17:25
		data, _ := mw.authSR.GetClientVerification(&users.ClientID{ClientID: params.ClientID})
		if !data.Data.EmailConfirm {
			cErr := cErrors.ResponseErrorModel{
				StandartCode: cErrors.StatusBadRequest,
				Message:      "confirm mail first",
			}
			return c.Status(400).JSON(cErr)
		}
		if data.Data.KycConfirm {
			cErr := cErrors.ResponseErrorModel{
				StandartCode: cErrors.StatusBadRequest,
				Message:      "you already have verified account",
			}
			return c.Status(400).JSON(cErr)
		}
		attempt := mw.usersUC.GetAuthKyc(&users.ClientID{ClientID: 1}) //было client.Data.ClientId
		t1 := attempt.Data.Auth.CreateTime
		t2 := utils.GetEuropeTime()
		delta := t2.Sub(t1).Hours()
		if delta < 24 {
			cErr := cErrors.ResponseErrorModel{
				StandartCode: cErrors.StatusBadRequest,
				Message:      "please wait 24 hours",
			}
			return c.Status(400).JSON(cErr)
		}
		return c.Next()
	}
}

func (mw *MDWManager) UpdateAdminLastEntry() fiber.Handler {
	return func(c *fiber.Ctx) error {
		//access := &admins.GetAdminByAccessTokenParams{Token: c.Get("Access")}
		//admin := mw.adminUC.GetAdminByAccessToken(access)
		//if admin.Error.InternalCode != 0 {
		//	return c.SendStatus(int(admin.Error.InternalCode))
		//}
		//params := &admins.UpdateAdminLastEntryParams{ID: admin.data.ClientID}
		//res := mw.adminUC.UpdateAdminLastEntry(params)
		//if res.Error.InternalCode != 0 {
		//	return c.SendStatus(int(admin.Error.InternalCode))
		//}
		return c.Next()
	}
}

func (mw *MDWManager) CheckAdminRole(params *admins.RoleParams) fiber.Handler {
	return func(c *fiber.Ctx) error {
		access := &admins.GetAdminByAccessTokenParams{Access: c.Get("Access")}
		admin := mw.adminUC.GetAdminByAccessToken(access)
		if admin.Error.InternalCode != 0 {
			return c.SendStatus(int(admin.Error.InternalCode))
		}
		if admin.Data.IsBlocked {
			return c.SendStatus(fiber.StatusUnauthorized)
		}
		adminRole := mw.adminUC.GetAdminRoles(&admins.GetAdminRoleParams{ClientID: admin.Data.ID})
		ok := utils.ContainsStr(adminRole.Data.Roles, params.Role)
		if !ok {
			if admin.Data.IsBlocked {
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}
		return c.Next()
	}
}

func (mw *MDWManager) CheckContentLength() fiber.Handler {
	return func(c *fiber.Ctx) error {
		length := c.Get("Content-Length")
		lengthInt := utils.StringToInt(length)
		if lengthInt > cConstants.MaxBodyLimit {
			return c.SendStatus(413)
		}
		return c.Next()
	}
}

func (mw *MDWManager) CheckContentLengthForImage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		length := c.Get("Content-Length")
		lengthInt := utils.StringToInt(length)
		if lengthInt > cConstants.MaxBodyLimitForAvatar {
			return c.SendStatus(413)
		}
		return c.Next()
	}
}
