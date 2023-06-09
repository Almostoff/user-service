package logger

import (
	"UsersService/config"
	"fmt"
	"os"
	"runtime"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Logger interface {
	InitLogger() error
	Debug(msg string)
	Debugf(template string, args ...interface{})
	Info(msg string)
	Infof(template string, args ...interface{})
	Warn(msg string)
	Warnf(template string, args ...interface{})
	Error(err error)
	Errorf(template string, args ...interface{})
	Fatal(msg string)
	Fatalf(template string, args ...interface{})
	Panic(msg string)
	Panicf(template string, args ...interface{})
}

// Logger
type ApiLogger struct {
	cfg    *config.Config
	tgBot  *tb.Bot
	logger zerolog.Logger
}

// App Logger constructor
func NewApiLogger(cfg *config.Config) *ApiLogger {
	return &ApiLogger{cfg: cfg}
}

var loggerLevelMap = map[string]zerolog.Level{
	"debug":    zerolog.DebugLevel,
	"info":     zerolog.InfoLevel,
	"warn":     zerolog.WarnLevel,
	"error":    zerolog.ErrorLevel,
	"panic":    zerolog.PanicLevel,
	"fatal":    zerolog.FatalLevel,
	"noLevel":  zerolog.NoLevel,
	"disabled": zerolog.Disabled,
}

func (a *ApiLogger) InitLogger() error {
	var w zerolog.LevelWriter
	a.logger = log.With().Caller().Logger()
	if a.cfg.Logger.InFile {
		logFile, err := os.OpenFile(a.cfg.Logger.FilePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		w = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile)
	} else {
		w = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{Out: os.Stdout})
	}

	if a.cfg.Logger.InTG {
		err := a.InitTG()
		if err != nil {
			return err
		}
		a.logger = zerolog.New(w).Level(loggerLevelMap[a.cfg.Logger.Level]).With().CallerWithSkipFrameCount(a.cfg.Logger.SkipFrameCount).Timestamp().Logger().Hook(a)
	} else {
		a.logger = zerolog.New(w).Level(loggerLevelMap[a.cfg.Logger.Level]).With().CallerWithSkipFrameCount(a.cfg.Logger.SkipFrameCount).Timestamp().Logger()
	}
	return nil
}

func (a *ApiLogger) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	if loggerLevelMap[a.cfg.Logger.Level] > level {
		return
	}
	go a.SendLogMessage(msg)
}

func (a *ApiLogger) Debug(msg string) {
	a.logger.Debug().Msg(msg)
}

func (a *ApiLogger) Debugf(template string, args ...interface{}) {
	a.logger.Debug().Msgf(template, args...)
}

func (a *ApiLogger) Info(msg string) {
	a.logger.Info().Msg(msg)
}

func (a *ApiLogger) Infof(template string, args ...interface{}) {
	a.logger.Info().Msgf(template, args...)
}

func (a *ApiLogger) Warn(msg string) {
	a.logger.Warn().Msg(msg)
}

func (a *ApiLogger) Warnf(template string, args ...interface{}) {
	a.logger.Warn().Msgf(template, args...)
}

func (a *ApiLogger) Error(err error) {
	a.logger.Error().Msg(err.Error())
}

func (a *ApiLogger) Errorf(template string, args ...interface{}) {
	fmt.Println("errorf")

	a.logger.Error().Msgf(template, args...)
}

func (a *ApiLogger) Panic(msg string) {
	a.logger.Panic().Msg(msg)
}

func (a *ApiLogger) Panicf(template string, args ...interface{}) {
	a.logger.Panic().Msgf(template, args...)
}

func (a *ApiLogger) Fatal(msg string) {
	a.logger.Fatal().Msg(msg)
}

func (a *ApiLogger) Fatalf(template string, args ...interface{}) {
	a.logger.Fatal().Msgf(template, args...)
}

func (a *ApiLogger) ErrorFull(error error) {
	_, fn, line, _ := runtime.Caller(1)
	msg := fmt.Sprintf("ERROR:\n%s :: %d :: %s", fn, line, error.Error())
	a.logger.Error().Stack().Err(error).Msg(msg)
}
