package handler

type Interface interface {
	// Ensure executes the handler specific business logic in order to complete
	// the given task, if possible.
	Ensure() error
}
