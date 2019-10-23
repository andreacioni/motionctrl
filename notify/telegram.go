package notify

import (
	"github.com/kpango/glg"

	"gopkg.in/telegram-bot-api.v4"
)

type TelegramNotifyService struct {
	apiToken string
	chatIds  []int64

	bot *tgbotapi.BotAPI
}

func (n *TelegramNotifyService) Authenticate() error {
	b, err := tgbotapi.NewBotAPI(n.apiToken)

	if err != nil {
		glg.Warnf("Failed to create new instance of bot API: %v", err)
	} else {
		glg.Debugf("Authorized on account %s", b.Self.UserName)

		n.bot = b
	}

	return nil
}

func (n *TelegramNotifyService) Notify(message, filename string) error {
	var err error

	if n.bot == nil {
		if err = n.Authenticate(); err != nil {
			return glg.Errorf("Can't authenticate to telegram bot API...")
		}
	}

	for _, chatID := range n.chatIds {

		if message != "" {
			msg := tgbotapi.NewMessage(chatID, message)

			if _, err = n.bot.Send(msg); err != nil {
				err = glg.Errorf("Failed to send notify message to %d: %v", chatID, err)
			}
		}

		if filename != "" {
			photo := tgbotapi.NewPhotoUpload(chatID, filename)

			if _, err = n.bot.Send(photo); err != nil {
				err = glg.Errorf("Failed to send notify photo to %d: %v", chatID, err)
			}
		}

	}

	return err
}

func (n *TelegramNotifyService) Stop() error {
	n.bot = nil
	n.apiToken = ""
	n.chatIds = nil

	return nil
}
