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
	"time"
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

type Response []struct {
	ID        string    `json:"id"`
	Text      []string  `json:"text"`
	CreatedAt time.Time `json:"created_at"`
}

type Carddeck struct {
	Name              string    `json:"name"`
	Code              string    `json:"code"`
	Description       string    `json:"description"`
	Unlisted          bool      `json:"unlisted"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	ExternalCopyright bool      `json:"external_copyright"`
	Category          string    `json:"category"`
	CallCount         string    `json:"call_count"`
	ResponseCount     string    `json:"response_count"`
	Author            struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	} `json:"author"`
	Rating string `json:"rating"`
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
		channel, _ := s.State.Channel(m.ChannelID)

		if invoked == "ping" {
			s.ChannelMessageSend(m.ChannelID, "Pong!")
		} else if invoked == "pong" {
			s.ChannelMessageSend(m.ChannelID, "Ping!")
		} else if invoked == "changeName" && m.Author.ID == "139386544275324928" {
			s.UserUpdate("", "", strings.Join(args, " "), s.State.User.Avatar, "")
			s.ChannelMessageSend(m.ChannelID, "Sucessfully changed name to: "+strings.Join(args, " "))
		} else if invoked == "cards" && channel.GuildID == "184394993375379457" {
			cors := strings.Join(args[1:], " ")
			if args[0] == "black" {
				url := "https://api.cardcastgame.com/v1/decks/23Z9M/calls"
				if strings.Count(cors, `_`) == 0 {
					s.ChannelMessageSend(m.ChannelID, "**Black cards** must contain a blank. `_`")
				} else {
					r, _ := regexp.Compile("_+")
					r2, _ := regexp.Compile("[.?!]$")
					cors = r.ReplaceAllString(cors, "_")
					if strings.Count(cors, "_") > 3 {
						s.ChannelMessageSend(m.ChannelID, "**Black cards** can only have a maximum of **3** blanks.")
					} else if !unicode.IsUpper(rune(cors[0])) && (cors[0] == '_' || !unicode.IsDigit(rune(cors[0]))) {
						s.ChannelMessageSend(m.ChannelID, "Are you sure the **Black card** begins with proper capitalization?")
					} else if !r2.MatchString(cors) {
						s.ChannelMessageSend(m.ChannelID, "Are you sure the **Black card** ends with proper punctiation?")
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
							s.ChannelMessageSend(m.ChannelID, "Successfully made `"+cors+"` a **Black card** for the HITMAN deck!")
						} else {
							defer res.Body.Close()
							body, _ := ioutil.ReadAll(res.Body)
							s.ChannelMessageSend(m.ChannelID, "Error, Code: **"+strconv.Itoa(res.StatusCode)+"**\n```json\n"+string(body)+"```")
						}
					}
				}
			} else if args[0] == "white" {
				r2, _ := regexp.Compile("[.]$")
				if r2.MatchString(cors) {
					s.ChannelMessageSend(m.ChannelID, "**White cards** shouldn't end with punctuation.")
				} else if strings.Contains(cors, "_") {
					s.ChannelMessageSend(m.ChannelID, "Are you *sure* you meant a **White card**? (If 100% sure, just tell Floretta#7311 to add it)")
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
						s.ChannelMessageSend(m.ChannelID, "Successfully made `"+cors+"` a **White card** for the HITMAN deck!")
					} else {
						defer res.Body.Close()
						body, _ := ioutil.ReadAll(res.Body)
						s.ChannelMessageSend(m.ChannelID, "Error, Code: **"+strconv.Itoa(res.StatusCode)+"**\n```json\n"+string(body)+"```")
					}
				}
			} else if args[0] == "view" {
				resp, _ := http.Get("https://api.cardcastgame.com/v1/decks/23Z9M/")
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				var r Carddeck
				err := json.Unmarshal(body, &r)
				if err != nil {
					fmt.Print("Error:", err)
				}
				i1, _ := strconv.Atoi(r.CallCount)
				i2, _ := strconv.Atoi(r.ResponseCount)
				total := strconv.Itoa(i1 + i2)
				s.ChannelMessageSend(m.ChannelID, "This is **"+r.Name+"** in its current state!\nThis deck is made by **"+r.Author.Username+"** and contains **"+total+"** cards! (B: "+r.CallCount+"/W: "+r.ResponseCount+")\n*<https://www.cardcastgame.com/browse/deck/"+r.Code+">*")
			} else if args[0] == "use" {
				s.ChannelMessageSend(m.ChannelID, "__How to use the **HITMAN** deck in CAH__\nMake a game on *<http://pyx-2.pretendyoure.xyz/zy/game.jsp>*\nIn the chat, use the command /addcardcast **23Z9M**\nProfit!")
			} else if args[0] == "search" {
				str := strings.Join(args[1:], " ")

				s.ChannelTyping(m.ChannelID)

				call, _ := http.Get("https://api.cardcastgame.com/v1/decks/23Z9M/calls")
				resp, _ := http.Get("https://api.cardcastgame.com/v1/decks/23Z9M/responses")

				defer resp.Body.Close()
				defer call.Body.Close()

				body1, _ := ioutil.ReadAll(call.Body)
				body2, _ := ioutil.ReadAll(resp.Body)

				var callresp Response
				var respresp Response

				err := json.Unmarshal(body1, &callresp)
				if err != nil {
					fmt.Print("Error:", err)
				}

				err = json.Unmarshal(body2, &respresp)
				if err != nil {
					fmt.Print("Error:", err)
				}

				send := "__**Black cards containing \"" + str + "\"**__\n```"
				for _, x := range callresp {
					if strings.Contains(strings.ToLower(strings.Join(x.Text, "_")), strings.ToLower(str)) {
						send += strings.Join(x.Text, "_") + "\n"
					}
				}
				send += "```__**White cards containing \"" + str + "\"**__\n```"
				for _, x := range respresp {
					if strings.Contains(strings.ToLower(strings.Trim(x.Text[0], " ")), strings.ToLower(str)) {
						send += strings.Trim(x.Text[0], " ") + "\n"
					}
				}
				send += "```"
				s.ChannelMessageSend(m.ChannelID, send)
			} else if args[0] == "help" {
				send := `__Help for the command: **cards**__
$cards white <content> - Adds a white card to the deck
$cards black <content> - Adds a black card to the deck
$cards view - See the current state of the deck
$cards use - Tells you how to use the deck
$cards search <string> - Searches in the cards for your string

__Example of adding a card__
$cards white Billy's genitals
$cards black Where did _ go?`
				s.ChannelMessageSend(m.ChannelID, send)
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
