package service

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofrs/uuid"
)

func list() *command {
	return &command{
		names: []string{"list", "ls"},
		handle: func(h *Handlers, e *messageCreatedEvent, args []string) error {
			reply := func(message string) error {
				if _, err := h.postMessage(e, message); err != nil {
					return fmt.Errorf("an error occurred while posting message: %w", err)
				}
				return nil
			}

			creatorID := uuid.FromStringOrNil(e.Message.User.ID)
			if creatorID == uuid.Nil {
				return reply("internal error: failed to get user id")
			}

			cs, err := h.repo.GetConverterByCreatorID(creatorID)
			if err != nil {
				log.Printf("An error occurred on GetConverterByCreatorID: %v\n", err)
				return reply("internal error: failed to get converters")
			}

			if !e.IsDM {
				if err := reply("DMを確認してください。"); err != nil {
					return err
				}
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("## Converters (%v)\n", len(cs)))
			sb.WriteString("\n")

			// 情報が増えたら追加する
			for _, c := range cs {
				sb.WriteString(fmt.Sprintf("### Converter `%s`\n", c.ID))
				sb.WriteString("\n")
				// TODO: チャンネル名
				sb.WriteString(fmt.Sprintf("- 投稿先チャンネル: !{\"type\":\"channel\",\"raw\":\"ココ\",\"id\":\"%s\"}", c.ChannelID))
			}

			// reply in DM
			_, err = h.postDirectMessage(creatorID.String(), sb.String())
			return err
		},
	}
}
