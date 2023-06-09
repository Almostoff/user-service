package loggerService

import (
	"UsersService/config"
	"UsersService/pkg/secure"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type logger struct {
	cfg *config.Config
}

func NewLoggerClient(cfg *config.Config) Logger {
	return &logger{
		cfg: cfg,
	}
}

func (l logger) DevLog(message string, levelId int64) {
	l.Log(&Log{
		Message:    message,
		ApiPublic:  l.cfg.LoggerService.DevPublic,
		ApiPrivate: l.cfg.LoggerService.DevPrivate,
		LevelId:    levelId,
	})
}

func (l logger) Log(log *Log) {
	body := LogBody{
		LevelId:   log.LevelId,
		ServiceId: l.cfg.LoggerService.ServiceId,
		Message:   log.Message,
	}
	bodyRaw, err := json.Marshal(body)
	if err != nil {
		fmt.Println("LOG ERROR", err)
		return
	}

	req, err := http.NewRequest("POST", l.cfg.LoggerService.Url+"logs", bytes.NewReader(bodyRaw))
	req.Header.Set("ApiPublic", log.ApiPublic)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Signature", secure.CalcSignature(log.ApiPrivate, string(bodyRaw)))

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("LOG ERROR", err)
	}
}
