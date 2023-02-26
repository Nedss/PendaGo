package main

import (
	"flag"
	"fmt"
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
	RoleId         = flag.String("r", "", "Role ID to add with discord command")
	TriggerCommand = flag.String("c", "", "Trigger command to add role on discord")
)

func init() { flag.Parse() }

func main() {

	discordBot, err := discordgo.New("Bot " + *BotToken)
	if err != nil {
		log.Fatal("error creating Discord session", err)
		return
	}

	discordBot.AddHandler(messageCreate)

	// Only listen receiving message events
	discordBot.Identify.Intents = discordgo.IntentsGuildMessages

	err = discordBot.Open()
	if err != nil {
		log.Fatal("error opening connection", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit or kill processus.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	discordBot.Close()
}

func createHelpEmbedMessage() *discordgo.MessageEmbed {
	helpEmbed := modules.NewEmbed().
		SetTitle("A l'aide Petit Penda !").
		SetDescription("Liste des commandes disponibles :").
		AddField("/poll", "ex : /poll 'Question' 'Réponse 1' 'Réponse 2' (9 réponses max)").
		AddField("/wowdiscord", "ex : /wowdiscord druid (classes: druid, monk, rogue, dh, hunter, shaman, priest, warlock, mage, warrior, paladin, dk, evoker)").
		AddField("/memes", "A vous de trouver.").
		SetColor(0x339A8C).MessageEmbed
	return helpEmbed
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	message := m.Content
	channelId := m.ChannelID

	if message == "/help" {
		messageEmbed := createHelpEmbedMessage()
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
	if *TriggerCommand == message {
		modules.AddWowRole(s, m, *RoleId)
	}
}
