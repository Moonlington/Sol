package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"strings"
)

var (
	prefix string
)

type Config struct {
	Token  string `json:"token"`
	Prefix string `json: "prefix"`
}

func main() {
	// Get the token from config.json
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Print("Error:", err)
	}
	var conf Config
	err = json.Unmarshal(content, &conf)
	if err != nil {
		fmt.Print("Error:", err)
	}

	// Create a new Discord session using the provided login information.
	// Use discordgo.New(Token) to just use a token for login.
	dg, err := discordgo.New("Bot " + conf.Token)
	prefix = conf.Prefix
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register messageCreate as a callback for the messageCreate events.
	dg.AddHandler(messageCreate)

	// Open the websocket and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	// Simple way to keep program running until CTRL-C is pressed.
	<-make(chan struct{})
	return
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if strings.HasPrefix(m.Content, prefix) {
		// Setting values for the commands
		args := strings.Split(m.Content[len(prefix):len(m.Content)], " ")
		invoked := args[0]
		args = args[1:len(args)]

		if invoked == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		}

		if invoked == "pong" {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		}

		if invoked == "changeName" && m.Author.ID == "139386544275324928" {
			s.UserUpdate("", "", strings.Join(args, " "), s.State.User.Avatar, "")
			s.ChannelMessageSend(m.ChannelID, "Sucessfully changed name to: "+strings.Join(args, " "))
		}
	}
}

// func messageCreate2(s *discordgo.Session, m *discordgo.MessageCreate) {

// 	_internal_channel := m.ChannelID
// 	_internal_author := m.Author

// 	if strings.HasPrefix(m.Content, dg.prefix) {
// 		args := strings.Split(m.Content[len(dg.prefix):len(m.Content)], " ")
// 		invoked := args[0]
// 		args = args[1:len(args)]

// 	}
// }
