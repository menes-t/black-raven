package service

import (
	"github.com/jasonlvhit/gocron"
	"github.com/menes-t/black-raven/config"
	"github.com/menes-t/black-raven/config/model"
	"github.com/menes-t/black-raven/logger"
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

	location, err := time.LoadLocation("Europe/Istanbul")
	now := time.Now().In(location)

	if err != nil {
		return
	}

	taskStartTime := time.Date(now.Year(), now.Month(), now.Day(), task.TaskConfig.StartingTimeAsHour, 0, 0, 0, now.Location())
	taskEndTime := time.Date(now.Year(), now.Month(), now.Day(), task.TaskConfig.EndingTimeAsHour, 0, 0, 0, now.Location())
	var startingDate time.Time

	var job *gocron.Job
	if now.Before(taskEndTime) && now.After(taskStartTime) {
		startingDate = now.Add(time.Second)
	} else if now.After(taskEndTime) {
		startingDate = time.Date(now.Year(), now.Month(), now.Day()+1, task.TaskConfig.StartingTimeAsHour, 0, 0, 0, now.Location())
	} else if now.Before(taskStartTime) {
		startingDate = time.Date(now.Year(), now.Month(), now.Day(), task.TaskConfig.StartingTimeAsHour, 0, 0, 0, now.Location())
	} else {
		//there is not any case here and the most important and only assumption is starting time < ending time
	}

	startingDate = postponeToPassWeekend(startingDate)

	job = taskScheduler.Every(task.TaskConfig.PeriodAsHour).Hours().From(&startingDate)
	logger.Logger().Info("Task is scheduled for " + startingDate.String())

	//TODO handle these cron job errors in some way
	job.Do(service.Remind, task)

	taskScheduler.Every(1).Day().At(strconv.Itoa(task.TaskConfig.EndingTimeAsHour) + ":00").Do(func() {
		taskScheduler.Remove(job)
		service.startNewTask(task, taskScheduler)
	})
}

func postponeToPassWeekend(now time.Time) time.Time {
	if now.Weekday() == time.Saturday {
		now = now.Add(2 * 24 * time.Hour)
	} else if now.Weekday() == time.Sunday {
		now = now.Add(1 * 24 * time.Hour)
	}
	return now
}
