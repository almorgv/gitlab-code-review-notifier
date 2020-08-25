package gitlabservice

import "github.com/xanzy/go-gitlab"

type FiringMergeRequest struct {
	MergeRequest      gitlab.MergeRequest
	FiringDiscussions []gitlab.Discussion
}
