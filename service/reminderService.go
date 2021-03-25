package service

import (
	"github.com/jasonlvhit/gocron"
	"github.com/menes-t/black-raven/config"
	"github.com/menes-t/black-raven/config/model"
	"github.com/menes-t/black-raven/model/response"
	"strconv"
	"time"
)

type Service struct {
	gitRepositoryServices map[string]GitRepositoryService
	messageServices       map[string]MessageService
	applicationConfig     config.ApplicationConfig
}

type MessageService interface {
	SendMessage(channelConfig model.MessageChannelConfig, mergeRequests []response.GitResponse)
}

type GitRepositoryService interface {
	Get(url string, token string) []response.GitResponse
}

type ReminderService interface {
	Remind(config model.Task)
	StartNewDay()
}

func NewReminderService(
	gitRepositoryServices map[string]GitRepositoryService,
	messageServices map[string]MessageService,
	applicationConfig config.ApplicationConfig,
) ReminderService {
	return &Service{
		gitRepositoryServices: gitRepositoryServices,
		messageServices:       messageServices,
		applicationConfig:     applicationConfig,
	}
}

func (service *Service) Remind(config model.Task) {
	mergeRequests := service.gitRepositoryServices["GitLab"].Get(config.GitRepositoryConfig.ApiUrl, config.GitRepositoryConfig.ApiToken)

	for _, channelConfig := range config.MessageConfig.Channels {
		service.messageServices["Slack"].SendMessage(channelConfig, mergeRequests) //TODO run this in a go routine
	}
}

func (service *Service) StartNewDay() {
	tasks := service.applicationConfig.Tasks
	taskScheduler := gocron.NewScheduler()
	for _, task := range tasks {
		now := time.Now()

		var job *gocron.Job
		if now.Hour() < task.TaskConfig.EndingTimeAsHour && now.Hour() > task.TaskConfig.StartingTimeAsHour {
			job = taskScheduler.Every(task.TaskConfig.PeriodAsHour).Hour().From(&now)
		} else if now.Hour() < task.TaskConfig.StartingTimeAsHour {
			job = taskScheduler.Every(task.TaskConfig.PeriodAsHour).Hour().At(service.applicationConfig.StartingTime)
		} else {
			job = taskScheduler.Every(task.TaskConfig.PeriodAsHour).Hour().At(strconv.Itoa(task.TaskConfig.StartingTimeAsHour) + ":00")
		}
		//TODO handle these cron job errors in some way
		job.Do(service.Remind, task)

		taskScheduler.Every(1).Day().At(strconv.Itoa(task.TaskConfig.EndingTimeAsHour) + ":00").Do(func() {
			taskScheduler.Clear()
			service.StartNewDay()
		})
	}
	<-taskScheduler.Start()
}
