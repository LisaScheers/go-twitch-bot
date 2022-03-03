package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

func TimerCallbackFn(msg twitch.PrivateMessage, client *twitch.Client) {
	// get args
	args := strings.Split(msg.Message[len("!timer "):], " ")
	if len(args) == 0 {
		client.Say(msg.Channel, "Usage: !timer <time>")
		return
	}
	// parse time
	timeParsed, err := time.ParseDuration(args[0])
	if err != nil {
		// debug
		log.Println(err)
		client.Say(msg.Channel, "Usage: !timer <time>")
		return
	}
	// start timer and show remaining time every 25% of the time
	timer := time.NewTimer(timeParsed)
	timerTime := time.Now().Add(timeParsed)
	for {
		select {
		case <-timer.C:
			client.Say(msg.Channel, "Timer done!")
			return
		case <-time.After(time.Minute):
			client.Say(msg.Channel, fmt.Sprintf("Timer remaining: %s", time.Until(timerTime)/time.Minute))
		}
	}

}
