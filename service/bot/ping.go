package bot

import "fmt"

func ping() *command {
	return &command{
		names: []string{"ping"},
		handle: func(h *Handlers, e *messageCreatedEvent, args []string) error {
			reply := func(message string) error {
				if _, err := h.postMessage(e, message); err != nil {
					return fmt.Errorf("an error occurred while posting message: %w", err)
				}
				return nil
			}

			return reply("Pong!")
		},
	}
}
