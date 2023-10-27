package router

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/sapphi-red/go-traq"
)

type postMessageArgs struct {
	channelID uuid.UUID
	message   string
}

// postMessage posts message to the specified channel.
func (h *Handlers) postMessage(channelID uuid.UUID, message string) (*traq.Message, error) {
	return h.throttle.Get(context.Background(), postMessageArgs{
		channelID: channelID,
		message:   message,
	})
}

func (h *Handlers) _postMessage(_ context.Context, args postMessageArgs) (*traq.Message, error) {
	m, _, err := h.api.ChannelApi.PostMessage(h.auth, args.channelID.String()).PostMessageRequest(traq.PostMessageRequest{
		Content: args.message,
	}).Execute()
	return m, err
}
