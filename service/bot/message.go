package bot

import (
	"log"
	"strings"

	"github.com/traPtitech/traq-ws-bot/payload"
)

type messageCreatedEvent struct {
	payload.Base
	Message payload.Message
	IsDM    bool
}

func (h *Handlers) MessageCreated(e *messageCreatedEvent) {
	// Do not process own message
	if e.Message.User.ID == h.botID.String() {
		return
	}
	// Do not process bot messages
	if e.Message.User.Bot {
		return
	}

	args := strings.Fields(e.Message.PlainText)

	for i, arg := range args {
		if strings.HasPrefix(arg, h.prefix) {
			cmdName := arg[len(h.prefix):]
			if c, ok := h.commands[cmdName]; ok {
				args[i] = cmdName
				// e.g. PlainText of "@BOT_example /ping arg1  arg2  " will be handed to command as
				// []string{"ping", "arg1", "arg2"}
				err := c.handle(h, e, args[i:])
				if err != nil {
					log.Printf("an error occurred while handling user command: %s\n", err)
				}
				return
			}
		}
	}
}
