package router

import (
	"github.com/antihax/optional"
	"github.com/gofrs/uuid"
	"github.com/sapphi-red/go-traq"
)

// postMessage posts message to the specified channel.
func (h *Handlers) postMessage(channelID uuid.UUID, message string) (*traq.Message, error) {
	m, _, err := h.api.ChannelApi.PostMessage(h.auth, channelID.String(), &traq.ChannelApiPostMessageOpts{
		PostMessageRequest: optional.NewInterface(traq.PostMessageRequest{
			Content: message,
			Embed:   false,
		}),
	})
	return &m, err
}
