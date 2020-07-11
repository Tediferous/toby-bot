package main

import (
        "github.com/bwmarrin/discordgo"
        log "github.com/sirupsen/logrus"
        "flag"
        "os"
        "os/signal"
        "syscall"
)

// Variables used for command line parameters
var (
	Token string
        Spam string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
        flag.StringVar(&Spam,"s","", "Spam Channel ID")
	flag.Parse()
}


func main(){
        log.Info("* Fight Song Plays *")
        dg, err := discordgo.New("Bot "+ Token)

        Check(err)

        dg.AddHandler(messageCreate)

        //Open socket
        Check(dg.Open())

        // Wait here until CTRL-C or other term signal is received.
	log.Info("Toby is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

        // :crown:
	if m.Content == "who's daddy" {
                s.ChannelMessageSend(Spam, "You are, King.")
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(Spam, "Pong!")
	}
}

func Check(err error){
        if err != nil {
                log.Error("Oops")
                log.Error(err)
        }
}

