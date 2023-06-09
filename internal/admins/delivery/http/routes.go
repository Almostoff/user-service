package http

import (
	"UsersService/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapAdminsRoutes(router fiber.Router, a *AdminHandler, mw *middleware.MDWManager) {
	routerGroup := router.Group("/admin")

	filegr := routerGroup.Group("/files", mw.CheckContentLengthForImage())
	filegr.Patch("/avatar", a.UpdateUserAvatar())

	bl := routerGroup.Group("/block")
	bl.Post("/", a.BlockUser())
	bl.Post("/unblock", a.UnBlockUser())
	bl.Post("/all", a.GetAllBlockUsers())
	kyc := routerGroup.Group("/kyc")
	kyc.Patch("/refresh", a.RefreshKyc())
	kyc.Post("/", a.GetAllClientKyc())
	change := routerGroup.Group("/change")
	change.Patch("/nickname", a.ChangeDndNickname())

}
