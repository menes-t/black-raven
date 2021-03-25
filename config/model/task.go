package model

type Task struct {
	GitRepositoryConfig GitRepositoryConfig
	MessageConfig       MessageConfig
	TaskConfig          TaskConfig
}

type GitRepositoryConfig struct {
	ApiUrl   string
	ApiToken string
}

type MessageConfig struct {
	ChannelNameWebHookUrlMap map[string]string
}

type TaskConfig struct {
	PeriodAsHour       uint64
	StartingTimeAsHour int
	EndingTimeAsHour   int
}
