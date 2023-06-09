package loggerService

import "UsersService/config"

var instance Logger = nil

func Init(cfg *config.Config) {
	instance = NewLoggerClient(cfg)
}

func GetInstance() Logger {
	return instance
}
