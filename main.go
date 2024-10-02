package main

import (
	"encoding/json"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
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
	Sesh.AddHandler(messageReactionAdd)
	Sesh.AddHandler(messageReactionRemove)

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
	trace(m)

	if rand.Intn(100) >= 97 {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🧢")
	}

	if rand.Intn(101) == 99 {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🔨")
	}

	// general debug and joke messages
	switch m.Content {
	case "go bears":
		s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" :bear:")
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
		if (m.Author.ID == GuildDaddy) && strings.Contains(params[0], "mute") {
			log.Debug("hammer? :eyes:")
			for _, u := range m.Mentions {
				if u.ID != s.State.User.ID {
					go mute(append([]string{u.ID, m.ChannelID}, params[1:]...)...)
				}
			}
			return
		} else if strings.Contains(params[0], "mute") {
			s.ChannelMessageSend(m.ChannelID, "you dont have the power to mute ,goofy")
			return
			} else if strings.Contains(params[0], "king me") {
				kingEm(s, m)
				return
			} else if strings.Contains(params[0], "beta me") {
				betaEm(s, m)
				return
		} else {

			s.MessageReactionAdd(m.ChannelID, m.ID, ":toby:732732965578211328")
		}

	}

}

func kingEm(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, _ := Sesh.GuildMember(Guild, m.Author.ID)
	Check(Sesh.GuildMemberRoleAdd(Guild, member.User.ID, "731317170386108509"))
	s.ChannelMessageSend(m.ChannelID, ":crown:")
}

func betaEm(s *discordgo.Session, m *discordgo.MessageCreate) {
	member, _ := Sesh.GuildMember(Guild, m.Author.ID)
	Check(Sesh.GuildMemberRoleRemove(Guild, member.User.ID, "731317170386108509"))
	s.ChannelMessageSend(m.ChannelID, ":b:")
}

func messageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	trace(r)

	if (r.Emoji.Name == "🔨") || (r.Emoji.Name == "nohammer") {
		tallyBanVotes(s, r.ChannelID, r.MessageID)
		return
	} else if (r.Emoji.Name == "lethalgun") && ( r.Member.User.ID == GuildDaddy){
		// Discord Owner has reacted with the kill word, ban them
		warrant :=  discernWhoToMute(r)
		go mute(warrant...)
	}

	return
}

func messageReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {

	trace(r)

	if (r.Emoji.Name == "🔨") || (r.Emoji.Name == "nohammer") {
		tallyBanVotes(s, r.ChannelID, r.MessageID)
		return
	}
	return
}

func discernWhoToMute(r *discordgo.MessageReactionAdd) []string {
	var warrant []string
	authorID :=  ""
	channel := ""
	sentence := "1h"
	// if this is a toby message
		// find who toby tagged
		// mute them
		
	warrant = append(warrant, authorID, channel, sentence)
	return warrant
}

func tallyBanVotes(s *discordgo.Session, channel, evidence string) {
	message, err := s.ChannelMessage(channel, evidence)
	Check(err)

	if isCaseClosed(message) {
		return
	} else {
		tally := make(map[string]int)
		tally["🔨"] = 0
		tally["nohammer"] = 0
		for _, reaction := range message.Reactions {
			tally[reaction.Emoji.Name] = reaction.Count
		}

		if tally["🔨"]-tally["nohammer"] >= 3 {
			// Check(s.ChannelMessageDelete(channel, evidence)) // delete message?
			daddy, err := s.User(GuildDaddy)
			Check(err)
			s.ChannelMessageSendReply(
				channel,
				"The people have decided that this message is shadow realm worthy. "+daddy.Mention(),
				message.Reference())
			closeCase(s, message)
			// go mute(message.Author.ID, channel, "1h")
		}
		return
	}
}

// react will react to the
// closeCase marks a message a seen for toby. This makes it so toby wont do judicial stuff with it anymore
func closeCase(s *discordgo.Session, message *discordgo.Message) error {
	s.MessageReactionAdd(message.ChannelID, message.ID, "🔒")
	return nil
}

// isCaseClosed checks to see if toby marked the case as closed
func isCaseClosed(message *discordgo.Message) bool {

	for _, r := range message.Reactions {
		if r.Emoji.Name == "🔒" && r.Me {
			return true
		}

	}

	return false
}

// messageDelete handles delete message events
func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	// TODO: make this function do something

	trace(m)
	// // Ignore all messages created by the bot itself
	// if m.Message.Author.ID == s.State.User.ID {
	// 	return
	// }

	// s.ChannelMessageSend(Spam, "@"+m.Message.Author.Username +" deleted the message:\n"+"`"+m.Message.Content+"`")
}

func mute(warrant ...string) {
	log.Debug("Getting out the hammer")
	trace(warrant)
	sceneOfTheCrime := warrant[1]
	// Dont ban the bot
	if warrant[0] == Sesh.State.User.ID {
		Sesh.ChannelMessageSend(sceneOfTheCrime, "I cant be banned sorry")
		return
	}

	member, err := Sesh.GuildMember(Guild, warrant[0])

	trace(member)
	Check(err)
	sentence := time.Hour

	if len(warrant) > 2 {
		sentence, err = time.ParseDuration(strings.TrimSpace(warrant[2]))
		if err != nil {
			sentence = 1 * time.Hour
			Sesh.ChannelMessageSend(sceneOfTheCrime, "I couldn't understand your time duration, so Ill set it to 1 hour")
		}
		maxSentence := 24 * time.Hour
		if sentence > maxSentence {
			sentence = maxSentence
			Sesh.ChannelMessageSend(sceneOfTheCrime, "the best I can do is 24 hours...")
		}
	}
	Sesh.ChannelMessageSend(sceneOfTheCrime, ":hammer:")

	//add banned role to member
	Check(Sesh.GuildMemberRoleAdd(Guild, member.User.ID, BanRole))

	//parse time and sleep
	time.Sleep(sentence)

	//remove banned role
	Check(Sesh.GuildMemberRoleRemove(Guild, member.User.ID, BanRole))

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
		log.Fatal(err)
	}
}

func trace(s interface{}) {
	res, _ := json.MarshalIndent(s, "", "  ")
	log.Debug(string(res))
}
