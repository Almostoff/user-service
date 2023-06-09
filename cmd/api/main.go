package main

import (
	"UsersService/config"
	"UsersService/internal/httpServer"
	"UsersService/pkg/httpErrorHandler"
	"UsersService/pkg/logger"
	"UsersService/pkg/loggerService"
	"UsersService/pkg/secure"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"log"
)

func main() {

	//var shield = secure.NewShield(os.Getenv("AE_KEY"))

	var shield = secure.NewShield("OEqBxziE9TdH18FLrRZ4Kr862z1Xh0UxgvgOSTU5OTc=")

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	//viperInstance, err := config.LoadConfig()
	//if err != nil {
	//	log.Fatalf("Cannot load config. Error: {%s}", err.Error())
	//}

	newConfig, err := config.NewConfig()

	//cfg, err := config.ParseConfig(viperInstance)
	//if err != nil {
	//	log.Fatalf("Cannot parse config. Error: {%s}", err.Error())
	//}
	//
	//config.DecryptConfig(cfg, shield)
	//loggerService.Init(cfg)

	config.DecryptConfig(newConfig, shield)
	loggerService.Init(newConfig)

	//appLogger := logger.NewApiLogger(cfg)
	//if err = appLogger.InitLogger(); err != nil {
	//	log.Fatalf("Cannot init logger. Error: {%s}", err.Error())
	//}

	appLogger := logger.NewApiLogger(newConfig)
	if err = appLogger.InitLogger(); err != nil {
		log.Fatalf("Cannot init logger. Error: {%s}", err.Error())
	}

	//appLogger.Infof("Logger [Level = %s; InFile = %t; filePath = %s; InTG: %t; chatID: %d]",
	//	cfg.Logger.Level, cfg.Logger.InFile, cfg.Logger.FilePath, cfg.Logger.InTG, cfg.Logger.ChatID,
	//)

	//errorHandler := httpErrorHandler.NewErrorHandler(cfg, appLogger)
	//s := httpServer.NewServer(cfg, appLogger, errorHandler, shield)
	//if err = s.Run(); err != nil {
	//	appLogger.ErrorFull(err)
	//}
	errorHandler := httpErrorHandler.NewErrorHandler(newConfig, appLogger)
	s := httpServer.NewServer(newConfig, appLogger, errorHandler, shield)
	if err = s.Run(); err != nil {
		appLogger.ErrorFull(err)
	}
}
