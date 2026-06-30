package telemetry

type Meter struct {
	service string
}

func NewMeter(service string) *Meter {
	return &Meter{service: service}
}

func (m *Meter) Counter(name string) Counter {
	return &noopCounter{}
}

func (m *Meter) Histogram(name string) Histogram {
	return &noopHistogram{}
}

type Counter interface {
	Add(value int64)
}

type Histogram interface {
	Record(value float64)
}

type noopCounter struct{}
func (c *noopCounter) Add(value int64) {}

type noopHistogram struct{}
func (h *noopHistogram) Record(value float64) {}
