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
		go service.messageServices["Slack"].SendMessage(channelConfig, mergeRequests)
	}
}

func (service *Service) StartNewDay() {
	tasks := service.applicationConfig.Tasks
	taskScheduler := gocron.NewScheduler()
	for _, task := range tasks {
		service.startNewTask(task, taskScheduler)
	}
	<-taskScheduler.Start()
}

func (service *Service) startNewTask(task model.Task, taskScheduler *gocron.Scheduler) {
	now := time.Now()

	var job *gocron.Job
	if now.Hour() < task.TaskConfig.EndingTimeAsHour && now.Hour() > task.TaskConfig.StartingTimeAsHour {
		now = now.Add(time.Second)
		job = taskScheduler.Every(task.TaskConfig.PeriodAsHour).Hour().From(&now)
	} else if (now.Hour() < task.TaskConfig.StartingTimeAsHour && (now.Hour() >= 0 || now.Hour() == 24)) || (now.Hour() >= task.TaskConfig.EndingTimeAsHour && now.Hour() < 24) {
		job = taskScheduler.Every(task.TaskConfig.PeriodAsHour).Hour().At(strconv.Itoa(task.TaskConfig.StartingTimeAsHour) + ":00")
	} else {
		//there is not any case here and the most important and only assumption is starting time < ending time
	}
	//TODO handle these cron job errors in some way
	job.Do(service.Remind, task)

	taskScheduler.Every(1).Day().At(strconv.Itoa(task.TaskConfig.EndingTimeAsHour) + ":00").Do(func() {
		taskScheduler.Remove(job)
		service.startNewTask(task, taskScheduler)
	})
}
