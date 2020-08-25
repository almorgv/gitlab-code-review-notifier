package firingservice

import (
	"time"

	"gitlab-code-review-notifier/pkg/log"
)

type FiringService struct {
	log.Loggable
}

func NewFiringService() *FiringService {
	return &FiringService{}
}

func (service *FiringService) ProcessAllConfigs(clients []*ConfiguredClient) {
	for _, config := range clients {
		service.ProcessConfig(config)
	}
}

func (service *FiringService) ProcessConfig(client *ConfiguredClient) {
	if len(client.Config.DiscussionFiringTimeout) > 0 {
		service.ProcessGroupMergeRequestDiscussions(client)
	}
	if len(client.Config.MergeRequestOldTimeout) > 0 {
		service.ProcessOldOpenedGroupMergeRequests(client)
	}
	if len(client.Config.MergeRequestReviewTimeout) > 0 {
		service.ProcessNeededReviewGroupMergeRequests(client)
	}
}

func (service *FiringService) ProcessOldOpenedGroupMergeRequests(client *ConfiguredClient) {
	service.Log().Infof("Start processing old opened merge requests in group %d", client.Config.GroupId)

	if len(client.Config.MergeRequestOldTimeout) == 0 {
		return
	}

	mrOldTimeout, err := time.ParseDuration(client.Config.MergeRequestOldTimeout)
	if err != nil {
		service.Log().Errorf("Failed to parse duration from %s: %v", client.Config.MergeRequestOldTimeout, err)
		return
	}

	oldMrs := client.Client.MergeRequests().GetOldOpenedGroupMergeRequests(client.Config.GroupId, mrOldTimeout)
	if len(oldMrs) > 0 {
		service.Log().Infof("Got %d old opened merge requests in group %d", len(oldMrs), client.Config.GroupId)
	}

	for _, mr := range oldMrs {
		client.Notifier.NotifyOldOpenedMergeRequest(mr, &client.Config)
	}

	service.Log().Infof("Finish processing old opened merge requests in group %d", client.Config.GroupId)
}

func (service *FiringService) ProcessNeededReviewGroupMergeRequests(client *ConfiguredClient) {
	service.Log().Infof("Start processing needed review merge requests in group %d", client.Config.GroupId)

	if len(client.Config.MergeRequestReviewTimeout) == 0 {
		return
	}

	mrReviewTimeout, err := time.ParseDuration(client.Config.MergeRequestReviewTimeout)
	if err != nil {
		service.Log().Errorf("Failed to parse duration from %s: %v", client.Config.MergeRequestReviewTimeout, err)
		return
	}

	mrs := client.Client.MergeRequests().GetNeededReviewGroupMergeRequests(client.Config.GroupId, mrReviewTimeout, client.Config.MergeRequestReviewersCount)
	if len(mrs) > 0 {
		service.Log().Infof("Got %d needed review merge requests in group %d", len(mrs), client.Config.GroupId)
	}

	for _, mr := range mrs {
		client.Notifier.NotifyNeededReviewMergeRequest(mr, &client.Config)
	}

	service.Log().Infof("Finish processing needed review merge requests in group %d", client.Config.GroupId)
}

func (service *FiringService) ProcessGroupMergeRequestDiscussions(client *ConfiguredClient) {
	service.Log().Infof("Start processing firing merge request discussions in group %d", client.Config.GroupId)

	if len(client.Config.DiscussionFiringTimeout) == 0 {
		return
	}

	discussionFiringTimeout, err := time.ParseDuration(client.Config.DiscussionFiringTimeout)
	if err != nil {
		service.Log().Errorf("Failed to parse duration from %s: %v", client.Config.DiscussionFiringTimeout, err)
		return
	}

	firingMergeRequests := client.Client.Discussions().GetFiringGroupMergeRequests(client.Config.GroupId, discussionFiringTimeout)
	if len(firingMergeRequests) > 0 {
		service.Log().Infof("Got %d firing merge request discussions in group %d", len(firingMergeRequests), client.Config.GroupId)
	}

	for _, fmr := range firingMergeRequests {
		client.Notifier.NotifyFiringMergeRequestDiscussions(fmr)
	}

	service.Log().Infof("Finish processing firing merge request discussions in group %d", client.Config.GroupId)
}
