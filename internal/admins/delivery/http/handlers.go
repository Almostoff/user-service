package http

import (
	"UsersService/config"
	"UsersService/internal/admins"
	"UsersService/internal/model"
	"UsersService/internal/users"
	"UsersService/pkg/cdnService"
	"UsersService/pkg/logger"
	"UsersService/pkg/loggerService"
	"UsersService/pkg/utils"
	"encoding/base64"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io/ioutil"
)

type AdminHandler struct {
	logger  *logger.ApiLogger
	cfg     *config.Config
	adminUC admins.AdminCase
	userUC  users.UseCase
	cdn     cdnService.UseCase
}

func NewAdminsHandlers(cfg *config.Config, logger *logger.ApiLogger, adminUC admins.AdminCase, cdn cdnService.UseCase, userUC users.UseCase) *AdminHandler {
	return &AdminHandler{cfg: cfg, logger: logger, adminUC: adminUC, userUC: userUC, cdn: cdn}
}

func (a AdminHandler) AdminSignIn() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params admins.AdminSignInParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := a.adminUC.AdminSignIn(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a AdminHandler) GetAdminRoles() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params1 admins.GetAdminByAccessTokenParams
		if err := utils.ReadRequest(ctx, &params1); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data1 := a.adminUC.GetAdminIDByNickname(&params1)
		if data1.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data1.Error.Message, data1.Error.InternalCode)
			// TODO: tg
		}

		data := a.adminUC.GetAdminRoles(&admins.GetAdminRoleParams{ClientID: data1.Data.ClientID})
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a AdminHandler) UnBlockUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params admins.UnBlockUserParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := a.adminUC.UnBlockUser(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a AdminHandler) BlockUser() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params admins.BlockUserParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		fmt.Println(params.BlockedUntil, params.ClientID)
		data := a.adminUC.BlockUser(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a AdminHandler) GetAllBlockUsers() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params admins.GetAllBlockUsersParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := a.adminUC.GetAllBlockUsers(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a AdminHandler) GetAllClientKyc() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params admins.Search
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := a.adminUC.GetAllClientKyc(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a AdminHandler) ChangeDndNickname() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params admins.ChangeDndNicknameParams
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := a.adminUC.ChangeDndNickname(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a AdminHandler) RefreshKyc() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params model.ClientID
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		data := a.adminUC.RefreshKyc(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			// TODO: tg
		}
		if data.Error.InternalCode == 0 {
			return ctx.JSON(data)
		}
		return ctx.Status(int(data.Error.StandartCode)).JSON(data)
	}
}

func (a *AdminHandler) UpdateUserAvatar() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var params users.UpdateUserAvatarParams
		file, err := ctx.FormFile("image")
		if err := utils.ReadRequest(ctx, &params); err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
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

		dataImage, err := a.cdn.SaveImage(&cdnService.SaveImageParams{ImageBase64: base64Encoding})
		if err != nil {
			loggerService.GetInstance().DevLog(fmt.Sprintf("Order -> SaveImg -> {%s}", err.Error()), 4)
			return ctx.JSON(fiber.StatusInternalServerError)

		}
		params.NewAvatar = dataImage.File
		clId := ctx.FormValue("client_id", "0")
		id := utils.StringToInt(clId)
		params.ClientID = id
		data := a.userUC.UpdateUserAvatar(&params)
		if data.Error.InternalCode != 0 {
			a.logger.Errorf("%s {%d}", data.Error.Message, data.Error.InternalCode)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("🙀Клиент %d не смог изменить аватар err: %s", params.ClientID, data.Error.Message), 1)
			go loggerService.GetInstance().DevLog(fmt.Sprintf("Users -> UpdateUserAvatar -> data = {%+v}, err = {%+v}", &params, &data), 4)
			return ctx.Status(int(data.Error.StandartCode)).JSON(data)
		}
		return ctx.JSON(data)
	}
}
