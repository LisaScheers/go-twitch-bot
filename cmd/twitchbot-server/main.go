package main

import (
	"github.com/LisaScheers/go-twitch-bot/pkg/bot"
	"github.com/gempir/go-twitch-irc/v3"
)

func main() {
	tb := bot.New(&bot.Config{
		Twitch: bot.Twitch{
			Username: "LisaScheers",
			Token:    "oauth:XXxXXXX",
			Channels: []string{"LisaScheers"},
		},
		Commands: []bot.Command{
			{
				Name:     "!timer",
				Func:     TimerCallbackFn,
				Prefix:   "!timer",
				Cooldown: 5,
			},
		},
		WhisperCommands: []bot.WhisperCommand{
			{
				Name: "!test",
				Func: func(msg twitch.WhisperMessage, client *twitch.Client) {
					client.Whisper(msg.User.Name, "test")
				},
				Prefix: "!",
			},
		},
	})
	err := tb.Setup()
	if err != nil {
		panic(err)
	}
}
