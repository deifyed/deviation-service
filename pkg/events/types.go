package events

import "time"

type Event struct {
	ID        string
	Timestamp time.Time
	Type      string
	Value     int64
}

type Client interface {
	Receive() Event
}
