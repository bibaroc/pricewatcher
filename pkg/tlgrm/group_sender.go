package tlgrm

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewGroupSender(token string, chatIDs ...int64) (func(string) error, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return func(text string) error {
		for i := range chatIDs {
			if _, err := bot.Send(tgbotapi.NewMessage(chatIDs[i], text)); err != nil {
				return err
			}
		}

		return nil
	}, nil
}
