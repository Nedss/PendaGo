package handlers

import (
	"pendago/modules"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageHandler(
	s *discordgo.Session,
	m *discordgo.MessageCreate,
	WowLogChanId string,
	WowChanId string,
	BotChanId string,
	TriggerCommand string,
	RoleId string,
	SWCRoleId string,
	SWCommand string,
) {
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
	if channelId == WowLogChanId {
		if strings.HasPrefix(message, "/analyselog") {
			formatLogLink, err := modules.GetLogAnalyzer(message)
			if err != nil {
				return
			}
			s.ChannelMessageSendReply(channelId, formatLogLink, m.Reference())
		}
	}

	// Case Wow channel
	if channelId == WowChanId {
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
	if channelId == BotChanId {
		if TriggerCommand == message {
			modules.AddWowRole(s, m, RoleId)
		}
		if SWCommand == message {
			modules.AddWowRole(s, m, SWCRoleId)
		}
	}
}
