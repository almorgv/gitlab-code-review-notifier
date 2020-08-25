package webhook

type Webhook interface {
	Send(text string) error
}
