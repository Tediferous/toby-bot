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
	Token      string
	Spam       string
	Sesh       *discordgo.Session
	Guild      string
	GuildDaddy string
	BanRole    string
)

func init() {
	log.SetLevel(log.DebugLevel)

	Token = os.Getenv("TOKEN")
	flag.StringVar(&Token, "t", Token, "Bot Token")
	flag.StringVar(&Spam, "s", "731514187335073823", "Spam Channel ID")
	flag.StringVar(&Guild, "g", "510456636263890949", "Guild ID")
	flag.StringVar(&BanRole, "b", "671763594643374100", "Ban Role ID")
	flag.Parse()
}

func main() {
	log.Info("* Fight Song Plays *")
	s, err := discordgo.New("Bot " + Token)
	Sesh = s

	Check(err)

	Sesh.AddHandler(messageCreate)
	Sesh.AddHandler(messageDelete)

	//Open socket
	Check(Sesh.Open())

	//Verify Guild Daddy
	g, err := Sesh.Guild(Guild)
	Check(err)
	GuildDaddy = g.OwnerID
	trace(GuildDaddy)

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Toby is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	Sesh.Close()
	log.Info("* Georgia on My Mind Plays *")
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
	case "go bears":
		s.ChannelMessageSend(m.ChannelID, ":bear:")
		return
	// :crown:
	case "who's alpha":
		if m.Author.ID == GuildDaddy {
			s.ChannelMessageSend(m.ChannelID, "You are, King.:crown:")
		} else {
			s.ChannelMessageSend(m.ChannelID, "Not you, :b:eta")
		}
		return
	}

	if isMentioned(s.State.User, m.Mentions) {
		log.Info("I heard my name")
		params := strings.Split(m.Content, ",")
		if m.Author.ID == GuildDaddy && strings.Contains(params[0], "ban") {
			log.Debug("hammer? :eyes:")
			for _, u := range m.Mentions {
				if u.ID != s.State.User.ID {
					go ban(append([]string{u.ID, m.ChannelID}, params[1:]...)...)
				}
			}
			return
			//TODO poll if message has '?'
		} else {
			s.ChannelMessageSend(m.ChannelID, ":b:etas dont have the power to ban")
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
	member, err := Sesh.GuildMember(Guild, warrant[0])
	sceneOfTheCrime := warrant[1]
	trace(member)
	Check(err)
	sentence := 5 * time.Hour

	if len(warrant) > 2 {
		sentence, err = time.ParseDuration(strings.TrimSpace(warrant[2]))
		if err != nil {
			sentence = 5 * time.Hour
		}
		maxSentence := 24 * time.Hour
		if sentence > maxSentence {
			sentence = maxSentence
			Sesh.ChannelMessageSend(sceneOfTheCrime, "the best I can do is 24 hours...")
		}
	}

	roles := member.Roles
	trace(roles)

	//remove all roles from member
	for _, role := range roles {
		Check(Sesh.GuildMemberRoleRemove(Guild, member.User.ID, role))
	}

	//add banned role to member
	Check(Sesh.GuildMemberRoleAdd(Guild, member.User.ID, BanRole))
	Sesh.ChannelMessageSend(sceneOfTheCrime, ":hammer:")

	//parse time and sleep
	time.Sleep(sentence)

	//remove banned role, add old roles
	for _, role := range roles {
		Check(Sesh.GuildMemberRoleAdd(Guild, member.User.ID, role))
	}
	Check(Sesh.GuildMemberRoleRemove(Guild, member.User.ID, BanRole))

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
		os.Exit(1)
	}
}

func trace(s interface{}) {
	res, _ := json.MarshalIndent(s, "", "  ")
	log.Debug(string(res))
}
