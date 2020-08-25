package firingservice

import (
	"gitlab-code-review-notifier/pkg/config"
	"gitlab-code-review-notifier/pkg/gitlabservice"
	"gitlab-code-review-notifier/pkg/notifier"
	"gitlab-code-review-notifier/pkg/webhook"
)

type ConfiguredClient struct {
	Client   *gitlabservice.Client
	Notifier *notifier.Notifier
	Config   config.FiringConfig
}

type ConfiguredClientFactory struct {
	gitlabClientFactory *gitlabservice.ClientFactory
	notifierFactory     *notifier.Factory
}

func NewConfiguredClientFactory(gitlabClientFactory *gitlabservice.ClientFactory, notifierFactory *notifier.Factory) *ConfiguredClientFactory {
	return &ConfiguredClientFactory{gitlabClientFactory: gitlabClientFactory, notifierFactory: notifierFactory}
}

func (f *ConfiguredClientFactory) MakeClient(config config.FiringConfig) (*ConfiguredClient, error) {
	gitlabClient, err := f.gitlabClientFactory.MakeClient(config.GitlabToken)
	if err != nil {
		return nil, err
	}
	return &ConfiguredClient{
		Client: gitlabClient,
		Notifier: f.notifierFactory.MakeWebhookNotifier(webhook.MattermostConfig{
			WebhookUrl:   config.WebhookUrl,
			Channel:      "",
			Username:     "",
			IconUrl:      "",
			DefaultColor: "#ff0000",
		}),
		Config: config,
	}, nil
}
