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
	gitRepositoryService GitRepositoryService
	messageService       MessageService
	applicationConfig    config.ApplicationConfig
}

type MessageService interface {
	SendMessage(host string, mergeRequests []response.GitResponse)
}

type GitRepositoryService interface {
	Get(url string, token string) []response.GitResponse
}

type ReminderService interface {
	Remind(config model.TaskConfig)
	StartNewDay()
}

func NewReminderService(
	gitRepositoryService GitRepositoryService,
	messageService MessageService,
	applicationConfig config.ApplicationConfig,
) ReminderService {
	return &Service{
		gitRepositoryService: gitRepositoryService,
		messageService:       messageService,
		applicationConfig:    applicationConfig,
	}
}

func (service *Service) Remind(config model.TaskConfig) {
	mergeRequests := service.gitRepositoryService.Get(config.ApiUrl, config.ApiToken)

	for _, webHookUrl := range config.ChannelNameWebHookUrlMap {
		service.messageService.SendMessage(webHookUrl, mergeRequests) //TODO run this in a go routine
	}
}

func (service *Service) StartNewDay() {
	tasks := service.applicationConfig.Tasks
	taskScheduler := gocron.NewScheduler()
	for _, task := range tasks {
		now := time.Now()

		var job *gocron.Job
		if now.Hour() < task.EndingTimeAsHour && now.Hour() > task.StartingTimeAsHour {
			job = taskScheduler.Every(task.PeriodAsHour).Hour().From(&now)
		} else if now.Hour() < task.StartingTimeAsHour {
			job = taskScheduler.Every(task.PeriodAsHour).Hour().At(service.applicationConfig.StartingTime)
		} else {
			job = taskScheduler.Every(task.PeriodAsHour).Hour().At(strconv.Itoa(task.StartingTimeAsHour) + ":00")
		}
		job.Do(service.Remind, task)

		taskScheduler.Every(1).Day().At(strconv.Itoa(task.EndingTimeAsHour) + ":00").Do(func() {
			taskScheduler.Clear()
			service.StartNewDay()
		})
	}
	<-taskScheduler.Start()
}
