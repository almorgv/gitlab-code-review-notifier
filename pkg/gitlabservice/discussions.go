package gitlabservice

import (
	"io/ioutil"
	"time"

	"github.com/xanzy/go-gitlab"

	"gitlab-code-review-notifier/pkg/log"
)

type DiscussionsService struct {
	client *gitlab.Client
	log.Loggable
}

func NewDiscussionsService(client *gitlab.Client) *DiscussionsService {
	return &DiscussionsService{client: client}
}

func (service *DiscussionsService) GetFiringGroupMergeRequests(groupId int, timeout time.Duration) []FiringMergeRequest {
	firingMergeRequests := make([]FiringMergeRequest, 0, 5)

	mrs := service.GetOpenedGroupMergeRequests(groupId)
	service.Log().Debugf("Received %d merge requests for group %d", len(mrs), groupId)

	for _, mr := range mrs {
		firingMergeRequestDiscussions := service.GetFiringMergeRequestDiscussions(mr, timeout)
		// MR is consider firing then it contains a firing discussions
		if len(firingMergeRequestDiscussions) > 0 {
			service.Log().Debugf(
				"Found %d firing discussions in merge request %d of project %d",
				len(firingMergeRequestDiscussions),
				mr.IID,
				mr.ProjectID,
			)
			firingMergeRequest := FiringMergeRequest{
				MergeRequest:      *mr,
				FiringDiscussions: firingMergeRequestDiscussions,
			}
			firingMergeRequests = append(firingMergeRequests, firingMergeRequest)
		}
	}

	return firingMergeRequests
}

func (service *DiscussionsService) GetOpenedGroupMergeRequests(groupId int) []*gitlab.MergeRequest {
	mrs, _, err := service.client.MergeRequests.ListGroupMergeRequests(groupId, &gitlab.ListGroupMergeRequestsOptions{
		State:       gitlab.String("opened"),
		ListOptions: gitlab.ListOptions{PerPage: 100},
	})
	if err != nil {
		service.Log().Warnf("Failed to ListGroupMergeRequests for group %d: %v", groupId, err)
		return nil
	}
	return mrs
}

func (service *DiscussionsService) GetFiringMergeRequestDiscussions(mr *gitlab.MergeRequest, timeout time.Duration) []gitlab.Discussion {
	outdatedDiscussions := make([]gitlab.Discussion, 0, 5)

	discussions := service.GetMergeRequestDiscussions(mr)

	service.Log().Debugf(
		"Got %d discussions for merge request %d in project %d",
		len(discussions),
		mr.IID,
		mr.ProjectID,
	)

	for _, discussion := range discussions {
		discussion.Notes = sanitizeNotes(discussion.Notes)
		if service.IsDiscussionFiring(mr, discussion, timeout) {
			service.Log().Debugf(
				"Found firing discussion %s in merge request %d of project %d",
				discussion.ID,
				mr.IID,
				mr.ProjectID,
			)
			outdatedDiscussions = append(outdatedDiscussions, *discussion)
		}
	}

	return outdatedDiscussions
}

func (service *DiscussionsService) GetMergeRequestDiscussions(mr *gitlab.MergeRequest) []*gitlab.Discussion {
	discussions := make([]*gitlab.Discussion, 0, 10)
	page := 1
	perPage := 100

	for {
		pageDiscussions := service.GetMergeRequestDiscussionsPaged(mr, page, perPage)
		discussions = append(discussions, pageDiscussions...)
		page++
		if len(pageDiscussions) < perPage {
			return discussions
		}
	}
}

func (service *DiscussionsService) GetMergeRequestDiscussionsPaged(mr *gitlab.MergeRequest, page int, perPage int) []*gitlab.Discussion {
	pageDiscussions, resp, err := service.client.Discussions.ListMergeRequestDiscussions(
		mr.ProjectID,
		mr.IID,
		&gitlab.ListMergeRequestDiscussionsOptions{
			Page:    page,
			PerPage: perPage,
		},
	)

	if err != nil {
		service.Log().Warnf("Failed to get discussions for merge request %d in project %d: %v", mr.IID, mr.ProjectID, err)
		return nil
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		service.Log().Warnf("Failed to get discussions for merge request %d in project %d: (%d) %s", mr.IID, mr.ProjectID, resp.StatusCode, body)
		return nil
	}

	return pageDiscussions
}

func (service *DiscussionsService) IsDiscussionFiring(mr *gitlab.MergeRequest, discussion *gitlab.Discussion, timeout time.Duration) bool {
	return isDiscussionResolvable(discussion) &&
		!isDiscussionResolved(discussion) &&
		isLastDiscussionParticipantAnAuthor(mr, discussion) &&
		isLastNoteOutdated(discussion, timeout)
}

func GetDiscussionParticipants(mr gitlab.MergeRequest, discussion gitlab.Discussion) []gitlab.BasicUser {
	participants := make(map[int]gitlab.BasicUser)

	// get only unique participants except an author of the merge request
	for _, note := range discussion.Notes {
		if note.Author.Username != mr.Author.Username {
			participants[note.Author.ID] = makeBasicUserFromNoteAuthor(note)
		}
	}

	// make the slice with participants
	res := make([]gitlab.BasicUser, 0, len(participants))
	for _, user := range participants {
		res = append(res, user)
	}

	return res
}

// Removes not resolvable notes in discussion e.g. system ones like "changed this line..."
func sanitizeNotes(notes []*gitlab.Note) []*gitlab.Note {
	res := make([]*gitlab.Note, 0, len(notes))
	for _, note := range notes {
		if note.Resolvable {
			res = append(res, note)
		}
	}
	return res
}

func isDiscussionResolvable(discussion *gitlab.Discussion) bool {
	return len(discussion.Notes) > 0 && discussion.Notes[0].Resolvable
}

func isDiscussionResolved(discussion *gitlab.Discussion) bool {
	return GetLastNoteInDiscussion(discussion).Resolved
}

func isLastDiscussionParticipantAnAuthor(mr *gitlab.MergeRequest, discussion *gitlab.Discussion) bool {
	return GetLastNoteInDiscussion(discussion).Author.Username == mr.Author.Username
}

func isLastNoteOutdated(discussion *gitlab.Discussion, timeout time.Duration) bool {
	// it is assumed that go-gitlab package returns timestamps in UTC
	return time.Now().UTC().After(GetLastNoteInDiscussion(discussion).CreatedAt.Add(timeout))
}

func GetLastNoteInDiscussion(discussion *gitlab.Discussion) *gitlab.Note {
	if len(discussion.Notes) == 0 {
		return nil
	}
	return discussion.Notes[len(discussion.Notes)-1]
}

func makeBasicUserFromNoteAuthor(note *gitlab.Note) gitlab.BasicUser {
	return gitlab.BasicUser{
		ID:        note.Author.ID,
		Username:  note.Author.Username,
		Name:      note.Author.Name,
		State:     note.Author.State,
		AvatarURL: note.Author.AvatarURL,
		WebURL:    note.Author.WebURL,
	}
}
