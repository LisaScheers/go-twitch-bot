// twitch bot
package bot

import (
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

// bot struct
type bot struct {
	client *twitch.Client
	config *Config
}

// config struct
type Config struct {
	Twitch struct {
		Username string
		Token    string
		Channel  string
	}
	Commands []struct {
		Name     string
		Func     func(msg twitch.PrivateMessage, client *twitch.Client)
		prefix   string
		cooldown int
	}
	WhisperCommands []struct {
		Name   string
		Func   func(msg twitch.WhisperMessage, client *twitch.Client)
		prefix string
	}
}

// new bot
func New(config *Config) *bot {
	// create new bot
	b := &bot{
		config: config,
	}
	// create new twitch client
	b.client = twitch.NewClient(config.Twitch.Username, config.Twitch.Token)
	// connect to twitch
	b.client.Connect()
	// handle messages
	b.client.OnPrivateMessage(b.handleMessage)
	// handle whispers
	b.client.OnWhisperMessage(b.handleWhisper)
	// return bot
	return b
}

// bot setup function
func (b *bot) setup() error {
	// setup client
	b.client = twitch.NewClient(b.config.Twitch.Username, b.config.Twitch.Token)
	b.client.OnConnect(func() {
		b.client.Join(b.config.Twitch.Channel)
	})
	b.client.OnPrivateMessage(b.handleMessage)
	b.client.OnWhisperMessage(b.handleWhisper)
	return b.client.Connect()
}

var cooldownMap map[string]map[string]time.Time = make(map[string]map[string]time.Time)

//handel message
func (b *bot) handleMessage(message twitch.PrivateMessage) {
	// handle comand
	for _, command := range b.config.Commands {
		if message.Message[0:len(command.prefix)] == command.prefix {
			//implement cooldown
			if _, ok := cooldownMap[command.Name]; !ok {
				cooldownMap[command.Name] = make(map[string]time.Time)
			}
			if _, ok := cooldownMap[command.Name][message.Channel]; !ok {
				cooldownMap[command.Name][message.Channel] = time.Now()
			}
			if time.Since(cooldownMap[command.Name][message.Channel]) < time.Duration(command.cooldown)*time.Second {
				continue
			}
			cooldownMap[command.Name][message.Channel] = time.Now()

			command.Func(message, b.client)
		}
	}
}

// handle whisper
func (b *bot) handleWhisper(message twitch.WhisperMessage) {
	// handle whisper
	for _, command := range b.config.WhisperCommands {
		if message.Message[0:len(command.prefix)] == command.prefix {
			command.Func(message, b.client)
		}
	}
}
