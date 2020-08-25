package webhook

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gitlab-code-review-notifier/pkg/log"
)

type MattermostConfig struct {
	WebhookUrl   string
	Channel      string
	Username     string
	IconUrl      string
	DefaultColor string
}

type Mattermost struct {
	config MattermostConfig
	log.Loggable
}

func NewMattermost(config MattermostConfig) *Mattermost {
	return &Mattermost{
		config: config,
	}
}

type MattermostMessage struct {
	Channel     string                 `json:"channel"`
	Username    string                 `json:"username"`
	IconUrl     string                 `json:"icon_url"`
	Attachments []MattermostAttachment `json:"attachments"`
}

type MattermostAttachment struct {
	Title     string `json:"title"`
	TitleLink string `json:"title_link"`
	Pretext   string `json:"pretext"`
	Text      string `json:"text"`
	Color     string `json:"color"`
}

func (m *Mattermost) Send(text string) error {
	return m.SendMessage(MattermostMessage{
		Channel:  m.config.Channel,
		Username: m.config.Username,
		IconUrl:  m.config.IconUrl,
		Attachments: []MattermostAttachment{
			{
				Color: m.config.DefaultColor,
				Text:  text,
			},
		},
	})
}

func (m *Mattermost) SendMessage(message MattermostMessage) error {
	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("serialize message %s to json: %v", message, err)
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{Transport: transport}

	resp, err := client.Post(m.config.WebhookUrl, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("do request: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("status code is %d body '%s'", resp.StatusCode, body[:500])
	}

	return nil
}
