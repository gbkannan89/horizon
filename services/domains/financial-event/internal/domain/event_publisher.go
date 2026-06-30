package domain

type DomainEvent interface {
	EventName() string
	EntityID() string
}

type Publisher interface {
	Publish(event DomainEvent) error
}
