package message

import (
	"fmt"
	"github.com/menes-t/black-raven/logger"
	"github.com/menes-t/black-raven/model/request"
	"github.com/menes-t/black-raven/model/response"
	"github.com/menes-t/black-raven/service"
	"github.com/menes-t/black-raven/service/message/http"
)

type SlackService struct {
	client http.Client
}

func NewSlackService(client http.Client) service.MessageService {
	return &SlackService{client: client}
}

func (service *SlackService) SendMessage(host string, mergeRequests []response.GitResponse) {
	blocks := []request.Block{
		{
			Type: "section",
			Text: &request.Markdown{
				Type: "mrkdwn",
				Text: "*Summary*",
			},
		},
		{
			Type: "divider",
		},
		{
			Type: "section",
			Fields: []request.Markdown{
				{
					Type: "mrkdwn",
					Text: "*Current Open Merge Requests*\n:alert: Number: $18,000 (ends in 53 days)\n:baklava: Earliest: $4,289.70\n:crown: Latest: $13,710.30",
				}, //TODO summary calculation
				{
					Type: "mrkdwn",
					Text: "*Today's Open Merge Requests*\n:alert-blue: Number: $18,000 (ends in 53 days)\n:baklava: Earliest: $4,289.70\n:crown: Latest: $13,710.30",
				}, //TODO summary calculation
			},
		},
		{
			Type: "section",
			Text: &request.Markdown{
				Type: "mrkdwn",
				Text: "*Merge Requests Awaiting Your Approval*",
			},
		},
		{
			Type: "divider",
		},
	}

	for _, mergeRequest := range mergeRequests {
		blocks = append(blocks, request.Block{
			Type: "section",
			Text: &request.Markdown{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*%s*\nSource Branch: *%s*\nTarget Branch: *%s*\nMerge Status: *%s*  \nCreated At: *%s*  \nUpdated At: *%s*  \nClick *<%s|me!>*",
					mergeRequest.Title,
					mergeRequest.SourceBranch,
					mergeRequest.TargetBranch,
					mergeRequest.MergeStatus,
					mergeRequest.CreatedAt,
					mergeRequest.UpdatedAt,
					mergeRequest.WebUrl,
				),
			},
		})
	}

	_, err := service.client.Get(host, request.SlackRequest{
		Type:   "home",
		Blocks: blocks,
	})

	if err != nil {
		logger.Logger().Error("Could not send message to slack")
	}
}
