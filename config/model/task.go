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
	Channels []MessageChannelConfig
}

type TaskConfig struct {
	PeriodAsHour       uint64
	StartingTimeAsHour int
	EndingTimeAsHour   int
}

type MessageChannelConfig struct {
	ChannelName          string
	WebHookUrl           string
	NotificationModifier string
	IconEmoji            string
	Username             string
}
