package bot

import (
	"github.com/sapphi-red/go-traq"
)

// getChannelPath gets the path to the channel.
func (h *Handlers) getChannelPath(channelID string) (string, error) {
	c, _, err := h.api.ChannelApi.GetChannel(h.auth, channelID).Execute()
	if err != nil {
		return "", err
	}
	if !c.ParentId.IsSet() || c.ParentId.Get() == nil {
		return c.Name, nil
	}
	p, err := h.getChannelPath(*c.ParentId.Get())
	if err != nil {
		return "", err
	}
	return p + "/" + c.Name, nil
}

// postMessage posts message to the channel in which the event happened.
func (h *Handlers) postMessage(e *messageCreatedEvent, message string) (*traq.Message, error) {
	embed := true
	m, _, err := h.api.ChannelApi.PostMessage(h.auth, e.Message.ChannelID).PostMessageRequest(traq.PostMessageRequest{
		Content: message,
		Embed:   &embed,
	}).Execute()
	return m, err
}

// postDirectMessage posts message to the specified user.
func (h *Handlers) postDirectMessage(userID string, message string) (*traq.Message, error) {
	embed := true
	m, _, err := h.api.UserApi.PostDirectMessage(h.auth, userID).PostMessageRequest(traq.PostMessageRequest{
		Content: message,
		Embed:   &embed,
	}).Execute()
	return m, err
}
