package notifier

import (
	"bytes"
	"fmt"
	"path"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/xanzy/go-gitlab"

	"gitlab-code-review-notifier/pkg/config"
	"gitlab-code-review-notifier/pkg/gitlabservice"
	"gitlab-code-review-notifier/pkg/log"
	"gitlab-code-review-notifier/pkg/webhook"
)

type Notifier struct {
	webhook          webhook.Webhook
	templatesBaseDir string
	log.Loggable
}

func NewNotifier(webhook webhook.Webhook, templatesBaseDir string) *Notifier {
	return &Notifier{
		webhook:          webhook,
		templatesBaseDir: templatesBaseDir,
	}
}

func (n *Notifier) NotifyOldOpenedMergeRequest(mr *gitlab.MergeRequest, config *config.FiringConfig) {
	templateFileName := "old_merge_request.gotpl"
	if err := n.notifyMessage(NewOldMergeRequestMessage(mr, config), templateFileName); err != nil {
		n.Log().Errorf("Failed to notify message: %v", err)
	}
}

func (n *Notifier) NotifyNeededReviewMergeRequest(mr *gitlabservice.MergeRequestWithParticipants, config *config.FiringConfig) {
	templateFileName := "needed_review_merge_request.gotpl"
	if err := n.notifyMessage(NewNeededReviewMergeRequestMessage(mr, config), templateFileName); err != nil {
		n.Log().Errorf("Failed to notify message: %v", err)
	}
}

func (n *Notifier) NotifyFiringMergeRequestDiscussions(fmr gitlabservice.FiringMergeRequest) {
	templateFileName := "firing_discussion.gotpl"

	for _, message := range MakeDiscussionMessages(fmr) {
		if err := n.notifyMessage(message, templateFileName); err != nil {
			n.Log().Errorf("Failed to notify message: %v", err)
		}
	}
}

func (n *Notifier) notifyMessage(data interface{}, templateFileName string) error {
	tplFilePath := path.Join(n.templatesBaseDir, templateFileName)
	tpl, err := template.New(templateFileName).
		Funcs(sprig.TxtFuncMap()).
		ParseFiles(tplFilePath)
	if err != nil {
		return fmt.Errorf("parse template %s: %v", tplFilePath, err)
	}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		return fmt.Errorf("compile message template: %v", err)
	}

	if err := n.webhook.Send(buf.String()); err != nil {
		return fmt.Errorf("send webhook message: %v", err)
	}

	return nil
}
