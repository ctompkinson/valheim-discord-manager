package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
)

var Token string

func main() {
	Token = os.Getenv("DISCORD_TOKEN")
	if Token == "" {
		log.Println("Token not found please set env variables DISCORD_TOKEN")
	}

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Println("error creating Discord session, ", err)
	}

	dg.AddHandler(getMessageHandler)

	err = dg.Open()
	if err != nil {
		log.Println("error creating Discord session, ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func getMessageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	found := false
	for _, role := range m.Member.Roles {
		if role == "809546453324267531" {
			found = true
		}
	}
	if !found {
		s.ChannelMessageSend(m.ChannelID, "You are not a godmin!")
		return
	}


	if !strings.Contains(m.Message.Content, "!valheim") {
		return
	}

	parts := strings.Split(m.Message.Content, " ")
	if len(parts) == 1 {
		s.ChannelMessageSend(m.ChannelID, "Must use start, stop or status")
		return
	}

	switch parts[1] {
	case "start":
		output, err := exec.Command("/home/vhserver/vhserver", "start").Output()
		log.Println(err)
		log.Println(string(output))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n%s\n", output, err))

	case "stop":
		output, err := exec.Command("/home/vhserver/vhserver", "stop").Output()
		log.Println(err)
		log.Println(string(output))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n%s\n", output, err))

	case "status":
		output, err := exec.Command("/home/vhserver/vhserver", "details").Output()
		log.Println(err)
		log.Println(string(output))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s\n%s\n", output, err))

	default:
		s.ChannelMessageSend(m.ChannelID, "Must use start, stop or status")
	}
	return
}
