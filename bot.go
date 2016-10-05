package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var (
	prefix string
	conf   Config
)

type Config struct {
	Token       string `json:"token"`
	Prefix      string `json: "prefix"`
	Hitmantoken string `json: "hitmantoken"`
	Gttoken     string `json: "gttoken"`
}

func main() {
	// Get the token from config.json
	content, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Print("Error:", err)
	}
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
		args = args[1:]

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

		if invoked == "cards" && m.ChannelID == "191760330517381121" {
			cors := strings.Join(args[1:], " ")
			if args[0] == "call" {
				url := "https://api.cardcastgame.com/v1/decks/23Z9M/calls"
				if strings.Count(cors, `_`) == 0 {
					s.ChannelMessageSend(m.ChannelID, "**Calls** must contain a blank. `_`")
				} else {
					r, _ := regexp.Compile("_+")
					r2, _ := regexp.Compile("[.?!]$")
					cors = r.ReplaceAllString(cors, "_")
					if strings.Count(cors, "_") > 3 {
						s.ChannelMessageSend(m.ChannelID, "**Calls** can only have a maximum of **3** blanks.")
					} else if !unicode.IsUpper(rune(cors[0])) && (cors[0] == '_' || !unicode.IsDigit(rune(cors[0]))) {
						s.ChannelMessageSend(m.ChannelID, "Are you sure the **Call** begins with proper capitalization?")
					} else if !r2.MatchString(cors) {
						s.ChannelMessageSend(m.ChannelID, "Are you sure the **Call** ends with proper punctiation?")
					} else {
						corsu := strings.Replace(cors, `_`, `",	"`, -1)
						payload := strings.NewReader(`{"calls":[{"text":["` + corsu + `"],"string":"` + cors + `","validation":{"state":"success","message":"Looks good!"},"eventChain":null}]}`)
						req, _ := http.NewRequest("POST", url, payload)

						req.Header.Add("content-type", "application/json")
						req.Header.Add("x-auth-token", conf.Hitmantoken)
						req.Header.Add("cache-control", "no-cache")
						req.Header.Add("postman-token", "1cc6b4e6-1baf-b98c-90c5-dfac8a222178")

						res, _ := http.DefaultClient.Do(req)
						if res.StatusCode == 201 {
							s.ChannelMessageSend(m.ChannelID, "Successfully made `"+cors+"` a **Call** for the HITMAN deck!")
						} else {
							defer res.Body.Close()
							body, _ := ioutil.ReadAll(res.Body)
							s.ChannelMessageSend(m.ChannelID, "Error, Code: **"+strconv.Itoa(res.StatusCode)+"**\n```json\n"+string(body)+"```")
						}
					}
				}
			} else if args[0] == "resp" {
				r2, _ := regexp.Compile("[.]$")
				if r2.MatchString(cors) {
					s.ChannelMessageSend(m.ChannelID, "**Responses** shouldn't end with punctuation.")
				} else {
					url := "https://api.cardcastgame.com/v1/decks/23Z9M/responses"
					payload := strings.NewReader(`{"responses":[{"text":["` + cors + `"],"string":"` + cors + `","validation":{"state":"success","message":"Looks good!"},"eventChain":null}]}`)
					req, _ := http.NewRequest("POST", url, payload)

					req.Header.Add("content-type", "application/json")
					req.Header.Add("x-auth-token", conf.Hitmantoken)
					req.Header.Add("cache-control", "no-cache")
					req.Header.Add("postman-token", "1cc6b4e6-1baf-b98c-90c5-dfac8a222178")

					res, _ := http.DefaultClient.Do(req)
					if res.StatusCode == 201 {
						s.ChannelMessageSend(m.ChannelID, "Successfully made `"+cors+"` a **Response** for the HITMAN deck!")
					} else {
						defer res.Body.Close()
						body, _ := ioutil.ReadAll(res.Body)
						s.ChannelMessageSend(m.ChannelID, "Error, Code: **"+strconv.Itoa(res.StatusCode)+"**\n```json\n"+string(body)+"```")
					}
				}
			} else if args[0] == "view" {
				s.ChannelMessageSend(m.ChannelID, "This is the **HITMAN** deck in its current state!\n*<https://www.cardcastgame.com/browse/deck/23Z9M>*")
			} else if args[0] == "use" {
				s.ChannelMessageSend(m.ChannelID, "__How to use the **HITMAN** deck in CAH__\nMake a game on *<http://pyx-2.pretendyoure.xyz/zy/game.jsp>*\nIn the chat, use the command /addcardcast **23Z9M**\nProfit!")
			}
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
