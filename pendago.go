package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"pendago/modules"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken       = flag.String("t", "", "Bot token")
	WowLogChanId   = flag.String("wl", "", "Channel ID to put WoW logs")
	WowChanId      = flag.String("wc", "", "Channel ID to put WoW commands")
	BotChanId      = flag.String("bc", "", "Chan of generic bot commands")
	RoleId         = flag.String("r", "", "Role ID to add with discord command")
	TriggerCommand = flag.String("c", "", "Trigger command to add role on discord")
	RoleBoost      = flag.String("rb", "", "Boost role ID")
	PendaRole      = flag.String("pr", "", "Penda Role ID")
	PendaGoldRole  = flag.String("pgr", "", "Penda Gold Role ID")
  SWCommand     = flag.String("sc", "", "Trigger command to add SWC role on discord")
  SWCRoleId     = flag.String("sr", "", "Role ID of SWC to add with discord command")
)

func init() { flag.Parse() }

func main() {
	discordBot, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatal("error creating Discord session", err)
		return
	}

	// discordBot.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	discordBot.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsGuildMembers

	err = discordBot.Open()
	if err != nil {
		log.Fatal("error opening connection", err)
		return
	}

	// Only listen receiving message events

	discordBot.AddHandler(messageCreate)
	discordBot.AddHandler(guildMemberBoost)

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit or kill processus.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discordBot.Close()
}

func guildMemberBoost(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	err := modules.BoostHandler(
		s,
		event,
		*RoleBoost,
		*PendaRole,
		*PendaGoldRole,
	)
	if err != nil {
		log.Fatal("Error during role attribution : ", err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	channelId := m.ChannelID

	if message == "/help" {
		messageEmbed := modules.CreateHelpEmbedMessage()
		s.ChannelMessageSendEmbed(channelId, messageEmbed)
	}

	// Case Wow log channel
	if channelId == *WowLogChanId {
		if strings.HasPrefix(message, "/analyselog") {
			formatLogLink, err := modules.GetLogAnalyzer(message)
			if err != nil {
				return
			}
			s.ChannelMessageSendReply(channelId, formatLogLink, m.Reference())
		}
	}

	// Case Wow channel
	if channelId == *WowChanId {
		// Wow Discord command
		if strings.HasPrefix(message, "/wowdiscord") {
			discordMessage := modules.GetDiscordClass(message)
			s.ChannelMessageSendReply(channelId, discordMessage, m.Reference())
		}
		// Memes
		if strings.HasPrefix(message, "/memes") {
			memeMessage, err := modules.GetMemes(message)
			if err != nil {
				return
			}
			s.ChannelMessageSendReply(channelId, memeMessage, m.Reference())
		}
	}
	// Global
	if strings.HasPrefix(message, "/poll") {
		modules.GeneratePollEmbedReaction(s, m)
	}
	if channelId == *BotChanId {
		if *TriggerCommand == message {
			modules.AddWowRole(s, m, *RoleId)
		}
    if *SWCommand == message {
			modules.AddWowRole(s, m, *SWCRoleId)
		}
	}
}
