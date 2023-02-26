package modules

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func getWowCommand(message string) (string, error) {
	splitMessage := strings.Split(message, " ")
	if len(splitMessage) != 2 {
		return "", fmt.Errorf("Invalid command. Can not split wow command.")
	}
	return splitMessage[1], nil

}

func GetDiscordClass(message string) string {

	class, err := getWowCommand(message)

	if err != nil {
		return "Commande incorrect.\n /help pour voir la liste des classes"
	}

	classLink := ""
	switch class {
	case "druid":
		classLink = "__Discord Druide :__\nhttps://discord.gg/0dWu0WkuetF87H9H"
	case "monk":
		classLink = "__Discord Moine :__\nhttp://discord.gg/peakofserenity"
	case "rogue":
		classLink = "__Discord Voleur :__\nhttps://discord.gg/0h08tydxoNhfDVZf"
	case "dh":
		classLink = "__Discord DH :__\nhttps://discord.gg/zGGkNGC"
	case "hunter":
		classLink = "__Discord Chasseur :__\nhttps://discord.gg/yqer4BX"
	case "shaman":
		classLink = "__Discord Chaman :__\nhttps://discord.gg/earthshrine"
	case "priest":
		classLink = "__Discord Prêtre :__\nhttps://discord.gg/HowToPriest"
	case "warlock":
		classLink = "__Discord Démoniste :__\nhttps://discord.gg/0onXDymd9Wpc2CEu"
	case "mage":
		classLink = "__Discord Mage :__\nhttps://discord.me/alteredtime"
	case "warrior":
		classLink = "__Discord Guerrier :__\nhttps://discord.gg/0pYY7932lTH4FHW6"
	case "paladin":
		classLink = "__Discord Paladin :__\nhttps://discord.gg/0dvRDgpa5xZHFfnD"
	case "dk":
		classLink = "__Discord DK :__\nhttps://discord.gg/acherus"
	case "evoker":
		classLink = "__Discord DK :__\nhttps://discord.gg/evoker"
	}

	if classLink == "" {
		return "Impossible de trouver le discord pour la classe demandée\n /help pour voir la liste des classes"
	}
	return classLink
}

func GetMemes(message string) (string, error) {

	meme, err := getWowCommand(message)
	if err != nil {
		return "", err
	}

	memeMessage := ""
	switch meme {
	case "sp":
		memeMessage = "**DEJA VU**\nhttps://www.youtube.com/watch?v=JQHtgq4G4zo"
	case "frostmage":
		memeMessage = "**UNLUCKY**\nhttps://www.youtube.com/watch?v=LocVPgHRhz8"
	}

	if memeMessage == "" {
		return "Pas de bonnes commandes, pas de memes !", fmt.Errorf("Command is not associated to a meme link")
	}
	return memeMessage, nil
}

func GetLogAnalyzer(message string) (string, error) {
	logLink, err := getWowCommand(message)
	if err != nil {
		return "", err
	}

	_, err = url.ParseRequestURI(logLink)
	if err != nil {
		return "", fmt.Errorf("Log url is not a valid url !")
	}

	urlId := path.Base(logLink)
	analyzerLink := "https://wowanalyzer.com/report/" + urlId

	return analyzerLink, nil
}

func AddWowRole(s *discordgo.Session, m *discordgo.MessageCreate, role string) {
	guildId := m.GuildID
	userId := m.Author.ID
	s.GuildMemberRoleAdd(guildId, userId, role)
	messageId := m.ID
	channelId := m.ChannelID
	s.ChannelMessageDelete(channelId, messageId)
}
