package service

import (
	"github.com/antihax/optional"
	"github.com/sapphi-red/go-traq"
)

// postMessage posts message to the channel in which the event happened.
func (h *Handlers) postMessage(e *messageCreatedEvent, message string) (*traq.Message, error) {
	m, _, err := h.api.ChannelApi.PostMessage(h.auth, e.Message.ChannelID, &traq.ChannelApiPostMessageOpts{
		PostMessageRequest: optional.NewInterface(traq.PostMessageRequest{
			Content: message,
			Embed:   true,
		}),
	})
	return &m, err
}
