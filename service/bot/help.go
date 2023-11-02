package bot

import (
	"fmt"
	"strings"
)

func help(commands []*command) *command {
	return &command{
		names: []string{"help"},
		handle: func(h *Handlers, e *messageCreatedEvent, args []string) error {
			var s strings.Builder
			s.WriteString("### Commands List\n")
			s.WriteString("Type each command without arguments to get help.\n")
			for _, cmd := range commands {
				for _, cmdName := range cmd.names {
					s.WriteString(fmt.Sprintf("- `/%v`\n", cmdName))
				}
			}
			_, err := h.postMessage(e, s.String())
			return err
		},
	}
}
