package toby

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
	log "github.com/sirupsen/logrus"
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
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

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
	// Sesh.AddHandler(messageReactionAdd)
	// Sesh.AddHandler(messageReactionRemove)

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
		// params := strings.Split(m.Content, ",")
		// if (m.Author.ID == GuildDaddy) && strings.Contains(params[0], "ban") {
		// 	log.Debug("hammer? :eyes:")
		// 	for _, u := range m.Mentions {
		// 		if u.ID != s.State.User.ID {
		// 			go ban(append([]string{u.ID, m.ChannelID}, params[1:]...)...)
		// 		}
		// 	}
		// 	return
		// } else if strings.Contains(params[0], "ban") {
		// 	s.ChannelMessageSend(m.ChannelID, ":b:etas dont have the power to ban")
		// 	return
		// 	// } else if strings.Contains(params[0], "king me") {
		// 	//         kingEm(s,m)
		// 	// return
		// 	// } else if strings.Contains(params[0], "beta me") {
		// 	//         betaEm(s,m)
		// 	// return
		// } else {

		s.MessageReactionAdd(m.ChannelID, m.ID, ":toby:732732965578211328")
		// }

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

func tallyBanVotes(s *discordgo.Session, channel, evidence string) {
	message, err := s.ChannelMessage(channel, evidence)
	Check(err)

	tally := make(map[string]int)
	tally["🔨"] = 0
	tally["nohammer"] = 0
	for _, reaction := range message.Reactions {
		tally[reaction.Emoji.Name] = reaction.Count
	}

	if tally["🔨"]-tally["nohammer"] >= 3 {
		Check(s.ChannelMessageDelete(channel, evidence))
		s.ChannelMessageSend(channel, "The people have spoken. Banning the perp and deleting the message so it can hurt us no more")
		go ban(message.Author.ID, channel, "2h")
	}
	return
}
func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	// TODO: make this function do something

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
	sceneOfTheCrime := warrant[1]
	// Dont ban the bot
	if warrant[0] == Sesh.State.User.ID {
		Sesh.ChannelMessageSend(sceneOfTheCrime, "I cant be banned sorry")
		return
	}

	member, err := Sesh.GuildMember(Guild, warrant[0])

	trace(member)
	Check(err)
	sentence := 2 * time.Hour

	if len(warrant) > 2 {
		sentence, err = time.ParseDuration(strings.TrimSpace(warrant[2]))
		if err != nil {
			sentence = 2 * time.Hour
			Sesh.ChannelMessageSend(sceneOfTheCrime, "I couldn't understand your time duration, so Ill set it to 2 hours")
		}
		maxSentence := 24 * time.Hour
		if sentence > maxSentence {
			sentence = maxSentence
			Sesh.ChannelMessageSend(sceneOfTheCrime, "the best I can do is 24 hours...")
		}
	}
	Sesh.ChannelMessageSend(sceneOfTheCrime, ":hammer:")

	roles := member.Roles
	trace(roles)

	//remove all roles from member
	for _, role := range roles {
		Check(Sesh.GuildMemberRoleRemove(Guild, member.User.ID, role))
	}

	//add banned role to member
	Check(Sesh.GuildMemberRoleAdd(Guild, member.User.ID, BanRole))

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
		log.Fatal(err)
	}
}

func trace(s interface{}) {
	res, _ := json.MarshalIndent(s, "", "  ")
	log.Debug(string(res))
}
