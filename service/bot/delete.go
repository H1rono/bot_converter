package bot

import (
	"fmt"

	"github.com/gofrs/uuid"

	"git.trap.jp/toki/bot_converter/repository"
)

func deleteConverter() *command {
	const help = "### Usage\n" +
		"\n" +
		"Converterを削除します。自身が持つconverterの一覧は、`/list`から確認することができます。\n" +
		"`/delete <converter_id>`\n" +
		"\n" +
		"- `/delete 00000000-0000-0000-0000-000000000000` 指定したconverterを削除します。"

	return &command{
		names: []string{"delete"},
		handle: func(h *Handlers, e *messageCreatedEvent, args []string) error {
			reply := func(message string) error {
				if _, err := h.postMessage(e, message); err != nil {
					return fmt.Errorf("an error occurred while posting message: %w", err)
				}
				return nil
			}

			if len(args) <= 1 {
				return reply(help)
			}

			creatorID := uuid.FromStringOrNil(e.Message.User.ID)
			if creatorID == uuid.Nil {
				return reply("internal error: failed to get user id")
			}

			converterID := uuid.FromStringOrNil(args[1])
			if converterID == uuid.Nil {
				return reply("Error: 正しいIDを指定してください。")
			}

			c, err := h.repo.GetConverter(converterID)
			if err != nil {
				if err == repository.ErrNotFound {
					return reply("Error: Converterが見つかりません。自身が所有権を持つconverter IDを指定してください。")
				} else {
					return reply("internal error: failed to get converter")
				}
			}
			if c.CreatorID != creatorID {
				return reply("Error: Converterが見つかりません。自身が所有権を持つconverter IDを指定してください。")
			}

			if err := h.repo.DeleteConverter(converterID); err != nil {
				return reply("internal error: failed to delete converter")
			}
			return reply("Converterを削除しました。")
		},
	}
}
