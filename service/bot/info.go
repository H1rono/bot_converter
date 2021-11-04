package bot

import (
	"fmt"

	"github.com/gofrs/uuid"

	"git.trap.jp/toki/bot_converter/repository"
)

func info() *command {
	const help = "### Usage\n" +
		"\n" +
		"Converterの情報を取得します。自身が持つconverterの一覧は、`/list`から確認することができます。\n" +
		"`/info <converter_id>`\n" +
		"\n" +
		"- `/info 00000000-0000-0000-0000-000000000000` 指定したconverterを削除します。"

	return &command{
		names: []string{"info"},
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

			// 情報が増えたら追加する
			return reply(fmt.Sprintf("Converter `%s`\n"+
				"\n"+
				// TODO: チャンネル名
				"- 投稿先チャンネル: !{\"type\":\"channel\",\"raw\":\"ココ\",\"id\":\"%s\"}",
				converterID, c.ChannelID))
		},
	}
}
