package notification

type Client interface {
	Send(title string, content string)
}
