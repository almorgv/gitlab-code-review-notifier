package config

import (
	"time"
)

type FiringConfig struct {
	Id                         int       `json:"id" db:"id"`
	GroupId                    int       `json:"group_id" db:"group_id"`
	GitlabToken                string    `json:"gitlab_token" db:"gitlab_token"`
	WebhookUrl                 string    `json:"webhook_url" db:"webhook_url"`
	DiscussionFiringTimeout    string    `json:"discussion_firing_timeout" db:"discussion_firing_timeout"`
	MergeRequestOldTimeout     string    `json:"merge_request_old_timeout" db:"merge_request_old_timeout"`
	MergeRequestOldMention     string    `json:"merge_request_old_mention" db:"merge_request_old_mention"`
	MergeRequestReviewTimeout  string    `json:"merge_request_review_timeout" db:"merge_request_review_timeout"`
	MergeRequestReviewersCount int       `json:"merge_request_reviewers_count" db:"merge_request_reviewers_count"`
	MergeRequestReviewMention  string    `json:"merge_request_review_mention" db:"merge_request_review_mention"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                  time.Time `json:"updated_at" db:"updated_at"`
}
