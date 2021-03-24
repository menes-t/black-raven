package model

type TaskConfig struct {
	ApiUrl                   string
	ApiToken                 string
	ChannelNameWebHookUrlMap map[string]string
	PeriodAsHour             uint64
	StartingTimeAsHour       int
	EndingTimeAsHour         int
}
