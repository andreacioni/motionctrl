package notify

import (
	"fmt"

	"github.com/kpango/glg"

	"gopkg.in/telegram-bot-api.v4"
)

type TelegramNotifyService struct {
	apiToken string
	chatIds  []string

	bot *tgbotapi.BotAPI
}

func (n *TelegramNotifyService) Authenticate() error {
	b, err := tgbotapi.NewBotAPI(n.apiToken)

	if err != nil {
		return fmt.Errorf("Failed to create new instance of bot API: %v", err)
	}

	glg.Debugf("Authorized on account %s", b.Self.UserName)

	n.bot = b

	return nil
}

func (n *TelegramNotifyService) Notify(filename string) error {
	return nil
}
