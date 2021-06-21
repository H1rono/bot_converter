package service

type command struct {
	names  []string
	handle func(h *Handlers, e *messageCreatedEvent, args []string) error
}

// commands returns the list of all commands.
func commands() []*command {
	return []*command{
		ping(),
		create(),
		info(),
		list(),
		deleteConverter(),
	}
}
