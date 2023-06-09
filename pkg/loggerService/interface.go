package loggerService

type Logger interface {
	Log(log *Log)
	DevLog(message string, levelId int64)
}
