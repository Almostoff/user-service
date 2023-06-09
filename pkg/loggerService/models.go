package loggerService

type Log struct {
	Message    string `json:"message"`
	ApiPublic  string `json:"-"`
	ApiPrivate string `json:"-"`
	LevelId    int64  `json:"level_id"`
}

type LogBody struct {
	LevelId   int64  `json:"level_id"`
	ServiceId int64  `json:"service_id"`
	Message   string `json:"message"`
}
