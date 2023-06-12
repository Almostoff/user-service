package http

import (
	"UsersService/internal/admins"
	"UsersService/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func MapUsersRoutes(router fiber.Router, h *UsersHandler, mw *middleware.MDWManager) {

	filegr := router.Group("/files", mw.CheckContentLengthForImage())
	filegr.Post("/avatar", h.UpdateUserAvatar())
	routerGroup := router.Group("/user", mw.CheckContentLength())
	routerGroup.Post("/sign_up", h.ClientSignUp()) //done
	routerGroup.Post("/sign_in", h.ClientSignIn()) //done
	routerGroup.Post("/sign_in_tg", h.ClientSignInTG())
	routerGroup.Get("/logout", h.Logout())                                             //done
	routerGroup.Get("/me/info", mw.UpdateUserLastActivity(), h.GetUserByAccessToken()) // done
	routerGroup.Post("/is_user_blocked", h.IsUserBlocked())                            // done

	routerGroup.Get("/get/rating", h.GetClientRating())                  // ctx - token - email - userID
	routerGroup.Post("/get_clients_statistics", h.GetClientStatistics()) // ctx - token - email - userID
	routerGroup.Post("/update_comment_client", h.UpdateCommentClient())  // ctx - token - email - userID
	routerGroup.Post("/update_comment", mw.CheckAdminRole(&admins.RoleParams{Role: "update_comment"}), h.UpdateCommentAdmin())

	strg := routerGroup.Group("/statistic")
	strg.Post("/add_comment/", h.AddComment())                  // ctx - token - email - userID
	strg.Get("/comment_exist/:internal_id", h.IsCommentExist()) // ctx - token - email - userID

	kycgr := router.Group("/kyc")
	kycgr.Post("/confirm", h.KycConfirm()) // ctx.Body - clientID
	kycusergr := routerGroup.Group("/kyc")
	kycusergr.Get("/init", mw.ValidateVerifications(), h.KycInit()) // ctx - token - email - clientid and uuid

	servgr := routerGroup.Group("/service")
	//mw.VerifySignatureMiddleware()) // в последнюю очередь
	servgr.Post("/add_notice", h.AddNotice()) // done
	servgr.Post("/verify/totp/with_id", h.VerifyTotpWithID())
	servgr.Get("/create_client", h.CreateClientDND())                      // nickname
	servgr.Get("/get_client_nickname/:client_id", h.GetUserNicknameByID()) // user id
	servgr.Get("/by_id/:user_id", h.GetUserByID())                         // req - user id
	servgr.Get("/get_user_id/by_access_token", h.GetUserIDByAccessToken()) // token - email - user id
	servgr.Get("/by_access_token", h.GetUserByAccessTokenWithID())
	servgr.Get("/by_nickname/:nickname", h.GetUserByNicknameWithID())

	statgr := servgr.Group("statistic")
	statgr.Post("/registration", h.GetRegistration())

	chrg := routerGroup.Group("/change",
		mw.UpdateUserLastActivity(),
		mw.ValidateToken(),
	)
	chrg.Post("/bio", h.UpdateUserBio())           // done
	chrg.Post("/nickname", h.UpdateUserNickName()) // использует сторонний сервис
	//chrg.Post("/avatar", h.UpdateUserAvatar())
	chrg.Post("/password", h.ChangePassword())
	chrg.Post("/language", h.ChangeLanguage()) // через таблицу lang
	chrg.Post("/tg", h.ChangeTg())

	rec := routerGroup.Group("/recover")
	rec.Post("/init", h.Recovery())
	rec.Post("/password/:hash", h.RecoveryPassword())

	conrg := routerGroup.Group("/confirm")
	conrg.Get("/email/:hash", h.ConfirmEmailByHash())
	conrg.Get("/email/code_init/withdraw", h.SendCodeToEmail("confirm_withdraw"))
	conrg.Post("/email/code/withdraw", h.CheckCodeFromEmail("confirm_withdraw"))
	conrg.Get("/init", h.ConfirmEmail())
	conrg.Post("/verify/totp", h.VerifyTotp())
	conrg.Post("/verify/totp_init", h.VerifyTotpInit())
	conrg.Get("/add/totp", h.AddTotp())

	tkgr := routerGroup.Group("/token")
	tkgr.Get("/refresh", h.RefreshAccess())
	tkgr.Get("/valid_access", h.ValidAccess())

	sesrg := routerGroup.Group("/sessions", mw.UpdateUserLastActivity(), mw.ValidateToken())
	sesrg.Get("/", h.GetUserSessions())
	sesrg.Post("/stop_session", h.StopSession())

	notrg := routerGroup.Group("/notice", mw.UpdateUserLastActivity()) // пустые?
	notrg.Get("/active", h.GetActiveNotice())
	notrg.Delete("/read/:internal_id", h.ReadNotice())

	routerGroup.Get("/info/all_reviews/:nickname", h.GetUserReviewsByNickname())
	routerGroup.Get("/:nickname", h.GetUserByNickname()) // done
}
