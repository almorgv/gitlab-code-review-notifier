package gitlabservice

import "github.com/xanzy/go-gitlab"

type MergeRequestWithParticipants struct {
	MergeRequest *gitlab.MergeRequest
	Participants []*gitlab.BasicUser
}
