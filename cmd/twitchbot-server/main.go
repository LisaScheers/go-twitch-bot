package main

import (
	"github.com/LisaScheers/go-twitch-bot/pkg/bot"
)

func main() {
	bot := bot.bot{
		config: &bot.Config{
			Twitch: &bot.Twitch{
				Username: "LisaScheers",
				Token:    "oauth:1cvh668q2689uglwcck08mv6v9oero",
				Channel:  "LisaScheers",
			},
			Commands: []bot.Command{
				{
					Name: "command",
					Func: func(msg bot.PrivateMessage, client *bot.Client) {
						client.Say(msg.Channel, "test back")
					},
					prefix:   "!test",
					cooldown: 10,
				},
			},
			WhisperCommands: []bot.Command{
				{
					Name:   "whisper",
					Func:   func(msg bot.WhisperMessage, client *bot.Client) {},
					prefix: "!",
				},
			},
		},
	}
	bot.setup()
}
