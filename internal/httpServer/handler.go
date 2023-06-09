package httpServer

import (
	adminsHttp "UsersService/internal/admins/delivery/http"
	adminsRepository "UsersService/internal/admins/repository"
	adminsUseCase "UsersService/internal/admins/usecase"
	"UsersService/internal/cConstants"
	"UsersService/internal/iConnection"
	iConnectionRepository "UsersService/internal/iConnection/repository"
	iConnectionUC "UsersService/internal/iConnection/usecase"
	"UsersService/internal/middleware"
	usersHttp "UsersService/internal/users/delivery/http"
	usersRepository "UsersService/internal/users/repository"
	usersUseCase "UsersService/internal/users/usecase"
	authService "UsersService/pkg/auth"
	"UsersService/pkg/cdnService"
	"UsersService/pkg/logger"
	ratingService "UsersService/pkg/rating"
	ssoService "UsersService/pkg/sso"
	"UsersService/pkg/storage"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	serverLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func (s *Server) MapHandlers(app *fiber.App, logger *logger.ApiLogger) error {

	db, err := storage.InitPsqlDB(s.cfg)
	if err != nil {
		return err
	}
	iConnectionRepo := iConnectionRepository.NewPostgresRepository(db, s.shield)
	iConnectionUC := iConnectionUC.NewIConnectionUsecase(iConnectionRepo)
	authConnection, err := iConnectionUC.GetInnerConnection(&iConnection.GetInnerConnectionParams{Name: &cConstants.AuthService})
	if err != nil {
		logger.ErrorFull(err)
		return err
	}
	ssoConnection, err := iConnectionUC.GetInnerConnection(&iConnection.GetInnerConnectionParams{Name: &cConstants.SsoService})
	if err != nil {
		logger.ErrorFull(err)
		return err
	}
	rateConnection, err := iConnectionUC.GetInnerConnection(&iConnection.GetInnerConnectionParams{Name: &cConstants.RatingService})
	//cdnConnection, err := iConnectionUC.GetInnerConnection(&iConnection.GetInnerConnectionParams{Name: &cConstants.CDN})
	if err != nil {
		logger.ErrorFull(err)
		return err
	}
	authSer := authService.GetClient(&authService.GetClientParams{
		Private: authConnection.Private,
		BaseUrl: authConnection.BaseUrl,
		Public:  authConnection.Public,
		Config:  s.cfg,
	})

	ssoSer := ssoService.GetClient(&ssoService.GetClientParams{
		Private: ssoConnection.Private,
		BaseUrl: ssoConnection.BaseUrl,
		Public:  ssoConnection.Public,
		Config:  s.cfg,
	})

	rateSer := ratingService.GetClient(&ratingService.GetClientParams{
		Private: rateConnection.Private,
		BaseUrl: rateConnection.BaseUrl,
		Public:  rateConnection.Public,
		Config:  s.cfg,
	})

	cdnSre := cdnService.GetClient(&cdnService.GetClientParams{
		IConn:  iConnectionUC,
		Config: s.cfg,
		//BaseUrl: *cdnConnection.BaseUrl,
		BaseUrl: "localhost:8282",
		//Public:  *cdnConnection.Public,
		//Private: *cdnConnection.Private,
		Public:  "test",
		Private: "test",
	})

	usersRepo := usersRepository.NewPostgresRepository(db, s.shield)
	usersUC := usersUseCase.NewUsersUsecase(logger, usersRepo, s.shield, authSer, ssoSer, rateSer, s.cfg.Kyc)
	usersHandlers := usersHttp.NewUsersHandlers(s.cfg, logger, usersUC, rateSer, authSer, ssoSer, cdnSre)

	adminsRepo := adminsRepository.NewPostgresRepository(db, s.shield)
	adminsUC := adminsUseCase.NewAdminUseCase(logger, adminsRepo, s.shield, authSer)
	adminsHandlers := adminsHttp.NewAdminsHandlers(s.cfg, logger, adminsUC, cdnSre, usersUC)

	app.Use(serverLogger.New())
	if _, ok := os.LookupEnv("LOCAL"); !ok {
		app.Use(recover.New())
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	mw := middleware.NewMDWManager(s.cfg, usersUC, adminsUC, authSer, iConnectionUC, *logger)
	usersHttp.MapUsersRoutes(app, usersHandlers, mw)
	adminsHttp.MapAdminsRoutes(app, adminsHandlers, mw)

	return nil
}
