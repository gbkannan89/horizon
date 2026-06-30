package events

import "context"

type Handler interface {
	Handle(ctx context.Context, envelope Envelope) error
}

type HandlerFunc func(ctx context.Context, envelope Envelope) error

func (f HandlerFunc) Handle(ctx context.Context, envelope Envelope) error {
	return f(ctx, envelope)
}
