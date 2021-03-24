package git

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/menes-t/black-raven/logger"
	"github.com/menes-t/black-raven/model/response"
	"github.com/menes-t/black-raven/service"
	"github.com/menes-t/black-raven/service/git/http"
)

type GitlabService struct {
	client http.Client
}

func NewGitlabService(client http.Client) service.GitRepositoryService {
	return &GitlabService{client: client}
}

func (service *GitlabService) Get(url string, token string) []response.GitResponse {
	responseAsByteArray, err := service.client.Get(url, token)

	if err != nil {
		logger.Logger().Error("Could not get merge requests")
		return []response.GitResponse{}
	}

	responseModel := make([]response.GitResponse, 0)
	err = jsoniter.Unmarshal(responseAsByteArray, &responseModel)

	if err != nil {
		logger.Logger().Error("Could not get merge requests")
		return []response.GitResponse{}
	}

	return responseModel
}
