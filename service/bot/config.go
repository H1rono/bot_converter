package bot

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gofrs/uuid"

	"git.trap.jp/toki/bot_converter/model"
	"git.trap.jp/toki/bot_converter/repository"
)

func config() *command {
	const help = "### Usage\n" +
		"\n" +
		"Converterの設定を編集します。\n" +
		"`/config <converter_id> [...args]`\n" +
		"\n" +
		"- `/config <converter_id>` 現在の設定を確認します。\n" +
		"- `/config <converter_id> pr-event-filter [filter1 [filter2]...]` 通知を行うPRのイベント型を複数の正規表現で指定します。\n" +
		"- `/config <converter_id> push-branch-filter [filter1 [filter2]...]` 通知を行うPushされたブランチ名を複数の正規表現で指定します。\n" +
		"\n" +
		"Examples:\n" +
		"- `/config <converter_id> pr-event-filter opened closed` PRが開かれた・閉じられた・マージされた場合にのみ通知します。\n" +
		"- `/config <converter_id> push-branch-filter main` mainブランチへのPushのみ通知します。\n"

	isValidRegexp := func(s string) bool {
		_, err := regexp.Compile(s)
		return err == nil
	}

	prEventFilterHandler := func(h *Handlers, e *messageCreatedEvent, args []string, conf *model.Config, reply func(string) error) error {
		filters := args[3:]

		// Validate
		for _, filter := range filters {
			if !isValidRegexp(filter) {
				return reply(fmt.Sprintf("invalid regexp filter: %v", filter))
			}
		}

		// Update
		conf.PREventTypesFilter = filters
		err := h.repo.SetConverterConfig(conf)
		if err != nil {
			return reply(fmt.Sprintf("internal error: failed to save config: %v", err))
		}
		return reply("Config saved!")
	}
	pushBranchFilterHandler := func(h *Handlers, e *messageCreatedEvent, args []string, conf *model.Config, reply func(string) error) error {
		filters := args[3:]

		// Validate
		for _, filter := range filters {
			if !isValidRegexp(filter) {
				return reply(fmt.Sprintf("invalid regexp filter: %v", filter))
			}
		}

		// Update
		conf.PushBranchFilter = filters
		err := h.repo.SetConverterConfig(conf)
		if err != nil {
			return reply(fmt.Sprintf("internal error: failed to save config: %v", err))
		}
		return reply("Config saved!")
	}

	return &command{
		names: []string{"config"},
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

			userID := uuid.FromStringOrNil(e.Message.User.ID)
			if userID == uuid.Nil {
				return reply("internal error: failed to get user id")
			}

			// retrieve converter
			converterID := uuid.FromStringOrNil(args[1])
			if converterID == uuid.Nil {
				return reply("Error: 正しいIDを指定してください。")
			}
			c, err := h.repo.GetConverter(converterID)
			if err != nil {
				if errors.Is(err, repository.ErrNotFound) {
					return reply("Error: Converterが見つかりません。自身が所有権を持つconverter IDを指定してください。")
				} else {
					return reply("internal error: failed to get converter")
				}
			}
			if c.CreatorID != userID {
				return reply("Error: Converterが見つかりません。自身が所有権を持つconverter IDを指定してください。")
			}

			// retrieve config
			conf, err := h.repo.GetConverterConfig(converterID)
			if err != nil && !errors.Is(err, repository.ErrNotFound) {
				return reply("internal error: failed to get config")
			}
			if errors.Is(err, repository.ErrNotFound) {
				conf = &model.Config{ConverterID: converterID}
			}

			if len(args) == 2 {
				return reply(fmt.Sprintf("Current config:\n"+
					"- pr-event-filter: %v\n"+
					"- push-branch-filter: %v\n",
					conf.PREventTypesFilter,
					conf.PushBranchFilter))
			}

			action := args[2]
			switch action {
			case "pr-event-filter":
				return prEventFilterHandler(h, e, args, conf, reply)
			case "push-branch-filter":
				return pushBranchFilterHandler(h, e, args, conf, reply)
			default:
				return reply(fmt.Sprintf("Unknown action: %v", action))
			}
		},
	}
}
