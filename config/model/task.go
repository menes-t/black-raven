package model

type TaskConfig struct {
	ChannelNameWebHookUrlMap map[string]string
	ApiUrl                   string
	ApiToken                 string
	PeriodAsHour             int
	StartingTimeAsHour       int
}
