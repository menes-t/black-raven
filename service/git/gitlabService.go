package git

import (
	"github.com/menes-t/black-raven/service"
	"github.com/menes-t/black-raven/service/git/http"
)

type GitlabService struct {
	client http.Client
}

func NewGitlabService(client http.Client) service.GitRepositoryService {
	return &GitlabService{client: client}
}

func (service *GitlabService) Get(url string) {
	//TODO gitlab logic here
}
