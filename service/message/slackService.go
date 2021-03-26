package message

import (
	"fmt"
	"github.com/menes-t/black-raven/config/model"
	"github.com/menes-t/black-raven/logger"
	"github.com/menes-t/black-raven/model/request"
	"github.com/menes-t/black-raven/model/response"
	"github.com/menes-t/black-raven/service"
	"github.com/menes-t/black-raven/service/message/http"
	"math"
	"time"
)

type SlackService struct {
	client http.Client
}

func NewSlackService(client http.Client) service.MessageService {
	return &SlackService{client: client}
}

func (service *SlackService) SendMessage(channelConfig model.MessageChannelConfig, mergeRequests []response.GitResponse) {
	var blocks []request.Block
	if len(mergeRequests) != 0 {
		//TODO implementing a block builder might be better (move all slack specific things to package slack)
		blocks = []request.Block{
			{
				Type: "header",
				Text: &request.Markdown{
					Type: "plain_text",
					Text: fmt.Sprintf("Hey %s! There are merge requests to look! :alert::alert:", channelConfig.NotificationModifier),
				},
			},
			{
				Type: "divider",
			},
		}
		summary := calculateSummary(mergeRequests)

		blocks = append(blocks, request.Block{
			Type: "section",
			Fields: []request.Markdown{
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Current Open Merge Requests*\n:alert: Number: %d\n:baklava: Earliest: %d hours\n:crown: Latest: %d hours",
						summary.MergeRequestCountTotal,
						summary.EarliestAsHour,
						summary.LatestAsHour,
					),
				},
				{
					Type: "mrkdwn",
					Text: fmt.Sprintf("*Today's Open Merge Requests*\n:alert-blue: Number: %d\n:baklava: Earliest: %d hours\n:crown: Latest: %d hours",
						summary.MergeRequestCountToday,
						summary.EarliestAsHourToday,
						summary.LatestAsHourToday,
					),
				},
			},
		})

		blocks = append(blocks, []request.Block{
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
		}...)

		fields := make([]request.Markdown, len(mergeRequests))
		for index, mergeRequest := range mergeRequests {
			fields[index] = request.Markdown{
				Type: "mrkdwn",
				Text: fmt.Sprintf("*<%s|%s>*\nAuthor: *%s*\nSource Branch: *%s*\nTarget Branch: *%s*\nMerge Status: *%s*  \nCreated At: *%s*  \nUpdated At: *%s*",
					mergeRequest.WebUrl,
					mergeRequest.Title,
					mergeRequest.Author.Name,
					mergeRequest.SourceBranch,
					mergeRequest.TargetBranch,
					mergeRequest.MergeStatus,
					mergeRequest.CreatedAt,
					mergeRequest.UpdatedAt,
				),
			}
		}

		blocks = append(blocks, request.Block{
			Type:   "section",
			Fields: fields,
		})
	} else {
		blocks = []request.Block{
			{
				Type: "header",
				Text: &request.Markdown{
					Type: "plain_text",
					Text: "*Hey congrats :omercan-party:! There is not any merge request!* :crown::crown::baklava::baklava:",
				},
			},
		}
	}

	_, err := service.client.Get(channelConfig.WebHookUrl, request.SlackRequest{
		Channel:   channelConfig.ChannelName,
		Type:      "home",
		Blocks:    blocks,
		IconEmoji: channelConfig.IconEmoji,
		Username:  channelConfig.Username,
	})

	if err != nil {
		logger.Logger().Error("Could not send message to slack")
	}
}

type Summary struct {
	MergeRequestCountTotal int
	MergeRequestCountToday int
	EarliestAsHour         int
	LatestAsHour           int
	EarliestAsHourToday    int
	LatestAsHourToday      int
}

//TODO this is slack related but not a concern of slack service (move all slack specific things to package slack)
func calculateSummary(mergeRequests []response.GitResponse) Summary {
	timeLayout := "2006-01-02T15:04:05.999Z07:00" //TODO time helper might be better
	summary := Summary{
		MergeRequestCountTotal: 0,
		MergeRequestCountToday: 0,
		EarliestAsHour:         math.MinInt32,
		LatestAsHour:           math.MaxInt32,
		EarliestAsHourToday:    math.MinInt32,
		LatestAsHourToday:      math.MaxInt32,
	}

	for _, mergeRequest := range mergeRequests {
		createdAt, err := time.Parse(timeLayout, mergeRequest.CreatedAt)

		if err != nil {
			continue
		}

		now := time.Now()
		elapsedTimeAsHoursAfterMergeRequest := int(now.Sub(createdAt) / (time.Hour / time.Nanosecond))

		if isToday(createdAt, now) {
			summary.MergeRequestCountToday += 1
			summary.MergeRequestCountTotal += 1

			if elapsedTimeAsHoursAfterMergeRequest > summary.EarliestAsHourToday {
				summary.EarliestAsHourToday = elapsedTimeAsHoursAfterMergeRequest
				summary.EarliestAsHour = elapsedTimeAsHoursAfterMergeRequest
			}

			if elapsedTimeAsHoursAfterMergeRequest < summary.LatestAsHourToday {
				summary.LatestAsHourToday = elapsedTimeAsHoursAfterMergeRequest
				summary.LatestAsHour = elapsedTimeAsHoursAfterMergeRequest
			}

		} else {
			summary.MergeRequestCountTotal += 1

			if elapsedTimeAsHoursAfterMergeRequest > summary.EarliestAsHour {
				summary.EarliestAsHour = elapsedTimeAsHoursAfterMergeRequest
			}

			if elapsedTimeAsHoursAfterMergeRequest < summary.LatestAsHour {
				summary.LatestAsHour = elapsedTimeAsHoursAfterMergeRequest
			}
		}
	}

	if summary.EarliestAsHour == math.MinInt32 {
		summary.EarliestAsHour = 0
	}

	if summary.LatestAsHour == math.MaxInt32 {
		summary.LatestAsHour = 0
	}

	if summary.EarliestAsHourToday == math.MinInt32 {
		summary.EarliestAsHourToday = 0
	}

	if summary.LatestAsHourToday == math.MaxInt32 {
		summary.LatestAsHourToday = 0
	}

	return summary
}

//TODO time helper might be better
func isToday(createdAt time.Time, now time.Time) bool {
	return createdAt.Day() == now.Day() && createdAt.Month() == now.Month() && createdAt.Year() == now.Year()
}
