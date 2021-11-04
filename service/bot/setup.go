package bot

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/sapphi-red/go-traq"
	traqbot "github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"

	"git.trap.jp/toki/bot_converter/repository"
)

type Handlers struct {
	repo     repository.Repository
	api      *traq.APIClient
	auth     context.Context
	botID    uuid.UUID
	prefix   string
	origin   string
	commands map[string]*command
}

// Start starts the bot service. Blocks on success.
func Start(c Config, repo repository.Repository) error {
	client := traq.NewAPIClient(traq.NewConfiguration())
	auth := context.WithValue(context.Background(), traq.ContextAccessToken, c.AccessToken)

	h := &Handlers{
		repo:     repo,
		api:      client,
		auth:     auth,
		botID:    c.BotID,
		prefix:   c.Prefix,
		origin:   c.Origin,
		commands: make(map[string]*command),
	}

	b, err := traqbot.NewBot(&traqbot.Options{
		AccessToken:   c.AccessToken,
		AutoReconnect: true,
	})
	if err != nil {
		return err
	}
	h.setUpHandlers(b)
	h.setUpCommands()

	return b.Start()
}

func (h *Handlers) setUpHandlers(b *traqbot.Bot) {
	b.OnMessageCreated(func(p *payload.MessageCreated) {
		h.MessageCreated(&messageCreatedEvent{
			Base:    p.Base,
			Message: p.Message,
			IsDM:    false,
		})
	})
	b.OnDirectMessageCreated(func(p *payload.DirectMessageCreated) {
		h.MessageCreated(&messageCreatedEvent{
			Base:    p.Base,
			Message: p.Message,
			IsDM:    true,
		})
	})
}

func (h *Handlers) setUpCommands() {
	cc := commands()
	for _, c := range cc {
		for _, name := range c.names {
			if _, ok := h.commands[name]; ok {
				log.Fatalf("command name conflict: " + name)
			}
			h.commands[name] = c
		}
	}
}
