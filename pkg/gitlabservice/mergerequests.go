package gitlabservice

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/xanzy/go-gitlab"

	"gitlab-code-review-notifier/pkg/log"
)

type MergeRequestsService struct {
	client *gitlab.Client
	log.Loggable
}

func NewMergeRequestsService(client *gitlab.Client) *MergeRequestsService {
	return &MergeRequestsService{client: client}
}

type predicate func(request *gitlab.MergeRequest) bool

func (r *MergeRequestsService) GetOldOpenedGroupMergeRequests(groupId int, timeout time.Duration) []*gitlab.MergeRequest {
	return r.filterGroupMergeRequests(groupId, func(mr *gitlab.MergeRequest) bool {
		return !mr.WorkInProgress && isMergeRequestNotUpdatedFor(mr, timeout)
	})
}

func (r *MergeRequestsService) GetNeededReviewGroupMergeRequests(groupId int, timeout time.Duration, reviewersCount int) []*MergeRequestWithParticipants {
	res := make([]*MergeRequestWithParticipants, 0)

	for _, mr := range r.GetOpenedGroupMergeRequests(groupId) {
		participants := r.GetMergeRequestsParticipants(mr)
		if !mr.WorkInProgress && isMergeRequestCreatedLongAgo(mr, timeout) && len(participants) < reviewersCount {
			res = append(res, &MergeRequestWithParticipants{
				MergeRequest: mr,
				Participants: r.GetMergeRequestsParticipants(mr),
			})
		}
	}

	return res
}

func (r *MergeRequestsService) GetOpenedGroupMergeRequests(groupId int) []*gitlab.MergeRequest {
	fullMrs := make([]*gitlab.MergeRequest, 0)

	mrs, resp, err := r.client.MergeRequests.ListGroupMergeRequests(groupId, &gitlab.ListGroupMergeRequestsOptions{
		State:       gitlab.String("opened"),
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		r.Log().Warnf("Failed to ListGroupMergeRequests for group %d: %v", groupId, err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		r.Log().Warnf("Failed to ListGroupMergeRequests for group %d: (%d) %s", groupId, resp.StatusCode, body)
		return nil
	}

	for _, mr := range mrs {
		fullMr, resp, err := r.client.MergeRequests.GetMergeRequestChanges(mr.ProjectID, mr.IID)
		if err != nil || resp.StatusCode != 200 {
			r.Log().Warnf("Failed to GetMergeRequestChanges for MR %d in project %d: %v", mr.IID, mr.ProjectID, err)
			return nil
		}
		fullMrs = append(fullMrs, fullMr)
	}

	return fullMrs
}

func (r *MergeRequestsService) GetMergeRequestsParticipants(mr *gitlab.MergeRequest) []*gitlab.BasicUser {
	allParticipants, resp, err := r.GetMergeRequestParticipants(mr.ProjectID, mr.IID)
	if err != nil {
		r.Log().Warnf("Failed to ListGroupMergeRequests for MR %d in project %d: %v", mr.IID, mr.ProjectID, err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		r.Log().Warnf("Failed to ListGroupMergeRequests for MR %d in project %d: (%d) %s", mr.IID, mr.ProjectID, resp.StatusCode, body)
		return nil
	}

	// removing the author from participants
	participants := make([]*gitlab.BasicUser, 0, len(allParticipants))
	for _, participant := range allParticipants {
		if participant.Username != mr.Author.Username {
			participants = append(participants, participant)
		}
	}

	return participants
}

// TODO remove after pull request will be approval
func (r *MergeRequestsService) GetMergeRequestParticipants(project int, mergeRequest int, options ...gitlab.RequestOptionFunc) ([]*gitlab.BasicUser, *gitlab.Response, error) {
	u := fmt.Sprintf("projects/%d/merge_requests/%d/participants", project, mergeRequest)

	req, err := r.client.NewRequest("GET", u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var p []*gitlab.BasicUser
	resp, err := r.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, err
}

func (r *MergeRequestsService) filterGroupMergeRequests(groupId int, predicate predicate) []*gitlab.MergeRequest {
	return r.filterMergeRequests(r.GetOpenedGroupMergeRequests(groupId), predicate)
}

func (r *MergeRequestsService) filterMergeRequests(mrs []*gitlab.MergeRequest, predicate predicate) []*gitlab.MergeRequest {
	result := make([]*gitlab.MergeRequest, 0)
	for _, mr := range mrs {
		if predicate(mr) {
			result = append(result, mr)
		}
	}
	return result
}

func isMergeRequestNotUpdatedFor(mr *gitlab.MergeRequest, duration time.Duration) bool {
	return isTimedOut(*mr.UpdatedAt, duration)
}

func isMergeRequestCreatedLongAgo(mr *gitlab.MergeRequest, timeout time.Duration) bool {
	return isTimedOut(*mr.CreatedAt, timeout)
}

func isTimedOut(t time.Time, timeout time.Duration) bool {
	// it is assumed that go-gitlab package returns timestamps in UTC
	return time.Now().UTC().After(t.Add(timeout))
}
