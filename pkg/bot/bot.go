// twitch bot
package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

// twitchBot struct
type twitchBot struct {
	client *twitch.Client
	config *Config
}

type Twitch struct {
	Username string
	Token    string
	Channels []string
}
type Command struct {
	Name     string
	Func     func(msg twitch.PrivateMessage, client *twitch.Client)
	Prefix   string
	Cooldown int
}
type WhisperCommand struct {
	Name   string
	Func   func(msg twitch.WhisperMessage, client *twitch.Client)
	Prefix string
}

// config struct
type Config struct {
	Twitch          Twitch
	Commands        []Command
	WhisperCommands []WhisperCommand
}

// new bot
func New(config *Config) *twitchBot {
	// create new bot
	b := &twitchBot{
		config: config,
	}
	return b
}

// bot setup function
func (b *twitchBot) Setup() error {
	// setup client
	b.client = twitch.NewClient(b.config.Twitch.Username, b.config.Twitch.Token)

	// join all channels
	for _, channel := range b.config.Twitch.Channels {
		b.client.Join(channel)
	}

	b.client.OnPrivateMessage(b.handleMessage)
	fmt.Println("setup commands complete")
	b.client.OnWhisperMessage(b.handleWhisper)
	fmt.Println("setup whisper commands complete")
	b.client.OnConnect(func() {
		fmt.Println("connected")
	})
	return b.client.Connect()
}

var cooldownMap map[string]map[string]time.Time = make(map[string]map[string]time.Time)

//handel message
func (b *twitchBot) handleMessage(message twitch.PrivateMessage) {
	// handle comand
	for _, command := range b.config.Commands {
		if message.Message[0:len(command.Prefix)] == command.Prefix {
			//implement cooldown
			if _, ok := cooldownMap[command.Name]; !ok {
				cooldownMap[command.Name] = make(map[string]time.Time)
			}
			if _, ok := cooldownMap[command.Name][message.Channel]; !ok {
				cooldownMap[command.Name][message.Channel], _ = time.Parse(time.RFC3339, "0001-01-01T00:00:00Z")
			}
			if time.Since(cooldownMap[command.Name][message.Channel]) < time.Duration(command.Cooldown)*time.Second {
				continue
			}
			cooldownMap[command.Name][message.Channel] = time.Now()
			log.Println("Running command:", command.Name)
			go command.Func(message, b.client)
		}
	}
}

// handle whisper
func (b *twitchBot) handleWhisper(message twitch.WhisperMessage) {
	// handle whisper
	for _, command := range b.config.WhisperCommands {
		if message.Message[0:len(command.Prefix)] == command.Prefix {
			go command.Func(message, b.client)
		}
	}
}

// register additonal commands
func (b *twitchBot) RegisterCommand(command Command) {
	b.config.Commands = append(b.config.Commands, command)
}
