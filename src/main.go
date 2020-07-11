package main

import (
        "github.com/bwmarrin/discordgo"
        log "github.com/sirupsen/logrus"
        "flag"
        "os"
        "os/signal"
        "syscall"
        "encoding/json"
        //"time"
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
        dg.AddHandler(messageDelete)

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
        trace(s)
        trace(m)

        // general debug and joke messages
        switch m.Content {
        case "ping":
                s.ChannelMessageSend(Spam, "Pong!")
                return
        // :crown:
        case "who's daddy":
                s.ChannelMessageSend(Spam, "You are, King.")
                return
        }

        // if m.Mentions

}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
        // TODO: make this function do something

        trace(s)
        trace(m)
	// // Ignore all messages created by the bot itself
	// if m.Message.Author.ID == s.State.User.ID {
	// 	return
	// }

        // s.ChannelMessageSend(Spam, "@"+m.Message.Author.Username +" deleted the message:\n"+"`"+m.Message.Content+"`")
}

func ban(warrant ...string){
        // MemberID, TimeInterval
        // member := warrant[0]
        // sentence := warrant[1]
        // roles := member.Roles

        //remove all roles from member

        //add banned role to member

        //parse time and sleep

        //remove banned role, add old roles
}

func poll(prompt string, quorum int){
        //send message to channel with prompt
        // options := map[string]int{
        //         ":thumbsup:": 1,
        //         ":thumbsdown:": 0,
        // }

        //watch message reactions
                //if either reaction reaches quorum
                //return options[reaction]
}

func Check(err error){
        if err != nil {
                log.Error("Oops")
                log.Error(err)
        }
}

func trace(s interface{}){
        res, _ := json.MarshalIndent(s,"","  ")
        log.Info(string(res))
}
