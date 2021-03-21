package message

import (
	"github.com/menes-t/black-raven/service"
	http2 "github.com/menes-t/black-raven/service/message/http"
)

type SlackService struct {
	client http2.Client
}

func NewSlackService(client http2.Client) service.MessageService {
	return &SlackService{client: client}
}

func (service *SlackService) SendMessage(host string, message string) {
	//TODO slack logic here
}
