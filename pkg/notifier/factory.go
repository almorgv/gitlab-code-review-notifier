package notifier

import "gitlab-code-review-notifier/pkg/webhook"

type Factory struct {
	templatesBaseDir string
}

func NewFactory(templatesBaseDir string) *Factory {
	return &Factory{templatesBaseDir: templatesBaseDir}
}

func (f Factory) MakeWebhookNotifier(config webhook.MattermostConfig) *Notifier {
	return NewNotifier(
		webhook.NewMattermost(config),
		"pkg/notifier/templates",
	)
}
