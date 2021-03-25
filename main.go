package main

import (
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

	messageServices := map[string]service.MessageService{"Slack": slackService}
	gitRepositoryServices := map[string]service.GitRepositoryService{"GitLab": gitlabService}

	reminderService := service.NewReminderService(gitRepositoryServices, messageServices, applicationConfig.GetConfig())

	reminderService.Remind(applicationConfig.GetConfig().Tasks[0]) //TODO delete this

	//reminderService.StartNewDay() TODO uncomment this
}
