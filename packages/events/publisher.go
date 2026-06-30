package events

import "context"

type Publisher interface {
	Publish(ctx context.Context, envelope Envelope) error
}
