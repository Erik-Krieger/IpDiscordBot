package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

const (
	CHANNEL_NAME = "ttt-server-channel"
)

// Variables used for command line parameters
var (
	Token string
)

type IP struct {
	Query string
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getIp() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(Ready)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func Ready(s *discordgo.Session, r *discordgo.Ready) {
	for _, g := range r.Guilds {
		ch, err := s.GuildChannels(g.ID)
		check(err)
		for _, c := range ch {
			if c.Name == CHANNEL_NAME {
				id := c.ID
				//s.ChannelMessageSend(id, getIp())
				c, err = s.ChannelEditComplex(id, &discordgo.ChannelEdit{
					Name:                 c.Name,
					Topic:                getIp(),
					NSFW:                 c.NSFW,
					Position:             c.Position,
					Bitrate:              c.Bitrate,
					UserLimit:            c.UserLimit,
					PermissionOverwrites: c.PermissionOverwrites,
					ParentID:             c.ParentID,
					RateLimitPerUser:     &c.RateLimitPerUser,
					Archived:             false,
					AutoArchiveDuration:  10080,
					Locked:               false,
					Invitable:            false,
				})
			}
		}
	}
	s.Close()
	os.Exit(0)
}
