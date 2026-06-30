package telemetry

type Span interface {
	End()
	SetAttributes(attrs map[string]string)
	RecordError(err error)
}

type noopSpan struct{}

func (s *noopSpan) End()                          {}
func (s *noopSpan) SetAttributes(attrs map[string]string) {}
func (s *noopSpan) RecordError(err error)         {}

type Tracer struct {
	service string
}

func NewTracer(service string) *Tracer {
	return &Tracer{service: service}
}

func (t *Tracer) StartSpan(name string) (Span, contextSetter) {
	return &noopSpan{}, func(ctx interface{}) interface{} { return ctx }
}

type contextSetter func(ctx interface{}) interface{}
