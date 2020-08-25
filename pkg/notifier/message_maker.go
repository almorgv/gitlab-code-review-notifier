package notifier

import (
	"time"

	"github.com/hako/durafmt"
	"github.com/xanzy/go-gitlab"

	"gitlab-code-review-notifier/pkg/config"
	"gitlab-code-review-notifier/pkg/gitlabservice"
)

type DiscussionMessage struct {
	MergeRequest  gitlab.MergeRequest
	Discussion    gitlab.Discussion
	Participants  []gitlab.BasicUser
	LastNote      gitlab.Note
	TimePassed    time.Duration
	TimePassedStr string
}

func NewDiscussionMessage(mergeRequest gitlab.MergeRequest, discussion gitlab.Discussion) DiscussionMessage {
	lastNote := *gitlabservice.GetLastNoteInDiscussion(&discussion)
	// it is assumed that go-gitlab package returns timestamps in UTC
	timePassed := time.Now().UTC().Sub(*lastNote.CreatedAt)
	return DiscussionMessage{
		MergeRequest:  mergeRequest,
		Discussion:    discussion,
		Participants:  gitlabservice.GetDiscussionParticipants(mergeRequest, discussion),
		LastNote:      lastNote,
		TimePassed:    timePassed,
		TimePassedStr: durafmt.Parse(timePassed).LimitFirstN(2).String(),
	}
}

func MakeDiscussionMessages(fmr gitlabservice.FiringMergeRequest) []DiscussionMessage {
	messages := make([]DiscussionMessage, 0)
	for _, discussion := range fmr.FiringDiscussions {
		messages = append(messages, NewDiscussionMessage(fmr.MergeRequest, discussion))
	}
	return messages
}

type OldMergeRequestMessage struct {
	MergeRequest           *gitlab.MergeRequest
	MergeRequestOldMention string
	TimePassed             time.Duration
	TimeSinceCreatedStr    string
	TimeSinceUpdatedStr    string
}

func NewOldMergeRequestMessage(mergeRequest *gitlab.MergeRequest, config *config.FiringConfig) OldMergeRequestMessage {
	// it is assumed that go-gitlab package returns timestamps in UTC
	timeSinceCreated := time.Now().UTC().Sub(*mergeRequest.CreatedAt)
	timeSinceUpdated := time.Now().UTC().Sub(*mergeRequest.UpdatedAt)
	return OldMergeRequestMessage{
		MergeRequest:           mergeRequest,
		MergeRequestOldMention: config.MergeRequestOldMention,
		TimeSinceCreatedStr:    durafmt.Parse(timeSinceCreated).LimitFirstN(2).String(),
		TimeSinceUpdatedStr:    durafmt.Parse(timeSinceUpdated).LimitFirstN(2).String(),
	}
}

type NeededReviewMergeRequestMessage struct {
	MergeRequest              *gitlab.MergeRequest
	Participants              []*gitlab.BasicUser
	MergeRequestReviewMention string
	TimeSinceCreatedStr       string
	TimeSinceUpdatedStr       string
}

func NewNeededReviewMergeRequestMessage(mrp *gitlabservice.MergeRequestWithParticipants, config *config.FiringConfig) NeededReviewMergeRequestMessage {
	// it is assumed that go-gitlab package returns timestamps in UTC
	timeSinceCreated := time.Now().UTC().Sub(*mrp.MergeRequest.CreatedAt)
	timeSinceUpdated := time.Now().UTC().Sub(*mrp.MergeRequest.UpdatedAt)
	return NeededReviewMergeRequestMessage{
		MergeRequest:              mrp.MergeRequest,
		Participants:              mrp.Participants,
		MergeRequestReviewMention: config.MergeRequestReviewMention,
		TimeSinceCreatedStr:       durafmt.Parse(timeSinceCreated).LimitFirstN(2).String(),
		TimeSinceUpdatedStr:       durafmt.Parse(timeSinceUpdated).LimitFirstN(2).String(),
	}
}
