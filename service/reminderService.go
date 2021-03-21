package service

import (
	"github.com/menes-t/black-raven/config/model"
)

type Service struct {
	gitRepositoryService GitRepositoryService
	messageService       MessageService
}

type MessageService interface {
	SendMessage(host string, message string)
}

type GitRepositoryService interface {
	Get(url string)
}

type ReminderService interface {
	Remind(config model.TaskConfig)
}

func NewReminderService(gitRepositoryService GitRepositoryService, messageService MessageService) ReminderService {
	return &Service{gitRepositoryService: gitRepositoryService, messageService: messageService}
}

func (service *Service) Remind(config model.TaskConfig) {
	//TODO reminder logic here
}
