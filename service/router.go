package service

import (
	"context"
	"log"

	"github.com/gofrs/uuid"
	"github.com/labstack/echo/v4"
	"github.com/sapphi-red/go-traq"
	traqbot "github.com/traPtitech/traq-bot"

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

func SetUp(c Config, repo repository.Repository) echo.HandlerFunc {
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

	eh := traqbot.EventHandlers{}
	h.setUpHandlers(eh)
	h.setUpCommands()
	server := traqbot.NewBotServer(c.VerificationToken, eh)

	return echo.WrapHandler(server)
}

func (h *Handlers) setUpHandlers(eh traqbot.EventHandlers) {
	eh.SetMessageCreatedHandler(func(p *traqbot.MessageCreatedPayload) {
		h.MessageCreated(&messageCreatedEvent{
			BasePayload: p.BasePayload,
			Message:     p.Message,
			IsDM:        false,
		})
	})
	eh.SetDirectMessageCreatedHandler(func(p *traqbot.DirectMessageCreatedPayload) {
		h.MessageCreated(&messageCreatedEvent{
			BasePayload: p.BasePayload,
			Message:     p.Message,
			IsDM:        true,
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
