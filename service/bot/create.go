package bot

import (
	"fmt"
	"log"

	"github.com/gofrs/uuid"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func create() *command {
	const help = "### Usage\n" +
		"\n" +
		"Converterを作成します。\n" +
		"`/create <here|channel> [secret]`\n" +
		"\n" +
		"- `/create here` このチャンネルに投稿するconverterを作成します。\n" +
		"- `/create #path/to/channel` #path/to/channel に投稿するconverterを作成します。\n" +
		"- `/create #path/to/channel my_webhook_secret` シークレットを指定した、 #path/to/channel に投稿するconverterを作成します。\n" +
		"  - :juuyo: シークレットを指定する場合は、シークレットの漏洩を防ぐためpublicチャンネルではなく本BOTへのDMで作成することを推奨します。"

	findEmbed := func(e *messageCreatedEvent, raw, embedType string) (payload.EmbeddedInfo, bool) {
		for _, embed := range e.Message.Embedded {
			if embedType != embed.Type {
				continue
			}
			if raw == embed.Raw {
				return embed, true
			}
		}
		return payload.EmbeddedInfo{}, false
	}

	return &command{
		names: []string{"create"},
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
				return reply("internal error: failed to get creator id")
			}

			// retrieve channel id
			var channelID uuid.UUID
			if args[1] == "here" {
				channelID = uuid.FromStringOrNil(e.Message.ChannelID)
			} else {
				embed, ok := findEmbed(e, args[1], "channel")
				if !ok {
					return reply(fmt.Sprintf("input error: %s に対応するembedを見つけることができませんでした。", args[1]))
				}
				channelID = uuid.FromStringOrNil(embed.ID)
			}
			if channelID == uuid.Nil {
				return reply("internal error: failed to get channel id")
			}

			// has secret
			var secret string
			if len(args) >= 3 {
				secret = args[2]
			}

			// create
			c, err := h.repo.CreateConverter(creatorID, channelID, secret)
			if err != nil {
				log.Printf("An error occurred on CreateConverter: %v\n", err)
				return reply("internal error: failed to create converter")
			}

			if !e.IsDM {
				if err := reply("Converterを作成しました。詳しくはBOTからのDMを参照してください。"); err != nil {
					return err
				}
			}

			// reply in DM
			_, err = h.postDirectMessage(creatorID.String(), fmt.Sprintf("Converterを作成しました。\n"+
				"Webhookの宛先を以下のURLに設定してください。\n"+
				"シークレットを指定した場合は、Webhookのシークレットをその値に設定してください。\n"+
				"\n"+
				":juuyo: **以下のURL及びconverterのIDは公開しないでください。公開してしまった場合、自由にメッセージを投稿される可能性があります。**\n"+
				"\n"+
				"- Gitea: %s/converters/%s/gitea\n"+
				"- GitHub: %s/converters/%s/github\n"+
				"\n"+
				"ここに列挙されていないサービスへの対応を望む場合は、BOT製作者に連絡してください。\n"+
				"\n"+
				"また、このconverterを削除する場合は `/delete %s` を実行してください。",
				h.origin, c.ID,
				h.origin, c.ID,
				c.ID))
			return err
		},
	}
}
