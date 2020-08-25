package gitlabservice

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/xanzy/go-gitlab"

	"gitlab-code-review-notifier/pkg/log"
)

type Client struct {
	discussions   *DiscussionsService
	mergeRequests *MergeRequestsService
	log.Loggable
}

func NewClient(gitlabToken string, gitlabUrl string) (*Client, error) {
	httpTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	httpClient := &http.Client{Transport: httpTransport}

	client, err := gitlab.NewClient(gitlabToken, gitlab.WithBaseURL(gitlabUrl), gitlab.WithHTTPClient(httpClient))
	if err != nil {
		return nil, fmt.Errorf("create gitlab client: %v", err)
	}

	return &Client{
		discussions:   NewDiscussionsService(client),
		mergeRequests: NewMergeRequestsService(client),
	}, nil
}

func (client *Client) Discussions() *DiscussionsService {
	return client.discussions
}

func (client *Client) MergeRequests() *MergeRequestsService {
	return client.mergeRequests
}

type ClientFactory struct {
	gitlabUrl string
}

func NewInstancedClientFactory(gitlabUrl string) *ClientFactory {
	return &ClientFactory{gitlabUrl: gitlabUrl}
}

func (f *ClientFactory) MakeClient(gitlabToken string) (*Client, error) {
	return NewClient(gitlabToken, f.gitlabUrl)
}
