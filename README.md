# gitlab-code-review-notifier

Notify about stale merge requests and code review discussions to slack or mattermost channel

## Install

### Helm chart

```
helm install -n gitlab \
--set ingress.hosts[0].host=gitlab-code-review-notifier.cluster.local \
--set ingress.hosts[0].paths[0]="/" \
--set gitlabUrl=https://gitlab \
gitlab-code-review-notifier gitlab-code-review-notifier
```

## API
### GET /clients
Get all clients

### GET /clients/:id
Get client by ID

### POST /clients
Add new client

##### Request body
`Content-Type: application/json`
```json
{
  "group_id": 250,
  "gitlab_token": "<YOUR_TOKEN>",
  "webhook_url": "https://mm/hooks/<YOUR_HOOK>",
  "merge_request_old_timeout": "24h",
  "merge_request_old_mention": "@all",
  "merge_request_review_timeout": "4h",
  "merge_request_reviewers_count": 2,
  "merge_request_review_mention": "@all",
  "discussion_firing_timeout": "2h"
}
```
`group_id` - ID of the group in gitlab to check code review in.

`gitlab_token` - token with `read_api` privileges of a user that has access to the specified group in gitlab.

`webhook_url` - mattermost incoming webhook url to send notifications.

`merge_request_old_timeout` - if set enables notification about old opened merge requests without WIP status.
Value is the duration passed since the merge request last update time.
The supported format is "24h30m" which max unit is hours.

`merge_request_old_mention` - what mention to use in the old opened merge requests notification message.
This parameter will be taken into account only if `merge_request_old_timeout` was set and notifications about old merge requests is enabled.
If not set default value `@all` will be used.

`merge_request_review_timeout` - if set enables notification about opened not WIP merge requests which doesn't have enough number of reviewers.
Value is the duration since the merge request created time must have passed to start checking reviewers.
The supported format is "24h30m" which max unit is hours.

`merge_request_reviewers_count` - number of reviewers required in each merge request.
This parameter will be taken into account only if `merge_request_review_timeout` was set and notifications about lack of reviewers is enabled.
If not set default value `1` will be used.

`merge_request_review_mention` - what mention to use in the lack of reviewers notification message.
This parameter will be taken into account only if `merge_request_review_timeout` was set and notifications about lack of reviewers is enabled.
If not set default value `@all` will be used.

`discussion_firing_timeout` - if set enables notification about unresolved discussions where last comment from the author of MR was left without an answer from reviewers.
Value is the duration passed since the last author comment creation in the discussion.
The supported format is "24h30m" which max unit is hours.

### PUT /clients/:id
Update existing client

##### Request body
Request body must contain **all** necessary fields that should be set in existing client as it performs full replace and all missing fields will be filled with default values.

`Content-Type: application/json`
```json
{
  "group_id": 250,
  "gitlab_token": "<YOUR_TOKEN>",
  "webhook_url": "https://mm/hooks/<YOUR_HOOK>",
  "merge_request_old_timeout": "24h",
  "merge_request_old_mention": "@all",
  "merge_request_review_timeout": "4h",
  "merge_request_reviewers_count": 2,
  "merge_request_review_mention": "@all",
  "discussion_firing_timeout": "2h"
}
```

### DELETE /clients/:id
Delete existing client
