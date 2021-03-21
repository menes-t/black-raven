package main

import (
	"github.com/jasonlvhit/gocron"
	"github.com/menes-t/black-raven/config"
	"github.com/menes-t/black-raven/service"
	"github.com/menes-t/black-raven/service/git"
	githttp "github.com/menes-t/black-raven/service/git/http"
	"github.com/menes-t/black-raven/service/message"
	messagehttp "github.com/menes-t/black-raven/service/message/http"
)

//TODO refactor http clients

func main() {
	applicationConfig, err := config.NewApplicationConfigGetter("resources/config.yml")
	if err != nil {
		panic(err)
	}

	gitRepositoryHTTPClient := githttp.NewGitRepositoryHTTPClient()
	gitlabService := git.NewGitlabService(gitRepositoryHTTPClient)

	messageHTTPClient := messagehttp.NewMessageHTTPClient()
	slackService := message.NewSlackService(messageHTTPClient)

	reminderService := service.NewReminderService(gitlabService, slackService)

	tasks := applicationConfig.GetConfig().Tasks

	for _, task := range tasks {
		taskScheduler := gocron.NewScheduler()
		taskScheduler.Every(1).Hour().Do(reminderService.Remind, task)
		<-taskScheduler.Start()
	}
}
