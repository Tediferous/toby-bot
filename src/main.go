package main

import (
	"encoding/json"
	"flag"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	Token   string
	Spam    string
	Sesh    discordgo.Session
	Guild   string
	BanRole string
)

func init() {
	log.SetLevel(log.DebugLevel)

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Spam, "s", "731514187335073823", "Spam Channel ID")
	flag.StringVar(&Guild, "g", "510456636263890949", "Guild ID")
	flag.StringVar(&BanRole, "b", "671763594643374100", "Ban Role ID")
	flag.Parse()
}

func main() {
	log.Info("* Fight Song Plays *")
	Sesh, err := discordgo.New("Bot " + Token)

	Check(err)

	Sesh.AddHandler(messageCreate)
	Sesh.AddHandler(messageDelete)

	//Open socket
	Check(Sesh.Open())

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Toby is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	Sesh.Close()
	log.Info("* Georgia Plays *")
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
	case "who's alpha":
		s.ChannelMessageSend(Spam, "You are, King.")
		return
	}

	if isMentioned(s.State.User, m.Mentions) {
		log.Info("I heard my name")
		params := strings.Split(m.Content, ",")
		if strings.Contains(params[0], "ban") { //TODO also check if m.Author has permission to ban
			log.Debug("hammer? :eyes:")
			for _, u := range m.Mentions {
				if u.ID != s.State.User.ID {
					go ban(append([]string{u.ID}, params[1:]...)...)
				}
			}
			return
			//TODO poll if message has '?'
		}
	}

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

func ban(warrant ...string) {
	log.Debug("Getting out the hammer")
        trace(warrant)
        Sesh.ChannelMessageSend(Spam,":hammer:")
	member, err := Sesh.GuildMember(Guild, warrant[0])
	Check(err)
	sentence, _ := time.ParseDuration("12h")

	if len(warrant) > 1 {
		sentence, err = time.ParseDuration(warrant[1])
		Check(err)
	}

	roles := member.Roles

	//remove all roles from member
	for _, role := range roles {
		Sesh.GuildMemberRoleRemove(Guild, member.User.ID, role)
	}

	//add banned role to member
	Sesh.GuildMemberRoleAdd(Guild, member.User.ID, BanRole)

	//parse time and sleep
	time.Sleep(sentence)

	//remove banned role, add old roles
	for _, role := range roles {
		Sesh.GuildMemberRoleAdd(Guild, member.User.ID, role)
	}
	Sesh.GuildMemberRoleRemove(Guild, member.User.ID, BanRole)

}

func poll(prompt string, quorum int) {
	//send message to channel with prompt
	// options := map[string]int{
	//         ":thumbsup:": 1,
	//         ":thumbsdown:": 0,
	// }

	//watch message reactions
	//if either reaction reaches quorum
	//return options[reaction]
}

func isMentioned(u *discordgo.User, m []*discordgo.User) bool {
	for _, user := range m {
		if u.ID == user.ID {
			return true
		}
	}
	return false
}
func Check(err error) {
	if err != nil {
		log.Error("Oops")
		log.Error(err)
	}
}

func trace(s interface{}) {
	res, _ := json.MarshalIndent(s, "", "  ")
	log.Debug(string(res))
}
