package handlers

import (
	"fmt"
	"pendago/modules"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func InteractionHandler(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	botChanID string,
	pendaRoleID string,
	rconHost string,
	rconPort string,
	rconPassword string,
) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()
	if data.Name != "minecraft" {
		return
	}

	if len(data.Options) == 0 || data.Options[0].Name != "whitelist_add" {
		return
	}

	if i.ChannelID != botChanID {
		respondInteraction(s, i, "Commande disponible uniquement dans le channel botcommand.")
		return
	}

	if !memberHasRole(i.Member, pendaRoleID) {
		respondInteraction(s, i, "Tu dois avoir le role Penda pour utiliser cette commande.")
		return
	}

	pseudo := ""
	for _, option := range data.Options[0].Options {
		if option.Name == "pseudo" {
			pseudo = option.StringValue()
			break
		}
	}

	if pseudo == "" {
		respondInteraction(s, i, "Pseudo manquant.")
		return
	}
	if !regexp.MustCompile("^[a-zA-Z0-9_]{2,16}$").MatchString(pseudo) {
		respondInteraction(s, i, "Le pseudo doit correspondre au format Minecraft. Alphanumérique et underscore seulement.")
		return
	}

	if err := deferInteraction(s, i); err != nil {
		respondInteraction(s, i, "Erreur lors du traitement de la commande.")
		return
	}

	output, err := modules.RunRCONCommand(
		rconHost,
		rconPort,
		rconPassword,
		"whitelist add "+pseudo,
	)
	lowerOutput := strings.ToLower(output)
	if strings.Contains(lowerOutput, "already") && strings.Contains(lowerOutput, "whitelist") {
		editInteractionResponse(s, i, fmt.Sprintf("%s est déjà présent dans la whitelist !", pseudo))
		return
	}
	if strings.Contains(lowerOutput, "does not exist") || strings.Contains(lowerOutput, "not found") {
		editInteractionResponse(s, i, fmt.Sprintf("%s n'existe pas.", pseudo))
		return
	}
	if err != nil {
		editInteractionResponse(s, i, "Erreur RCON.")
		return
	}

	editInteractionResponse(
		s,
		i,
		fmt.Sprintf("%s a été ajouté à la whitelist par %s", pseudo, i.Member.User.Username),
	)
}

func memberHasRole(member *discordgo.Member, roleID string) bool {
	if member == nil {
		return false
	}
	for _, role := range member.Roles {
		if role == roleID {
			return true
		}
	}
	return false
}

func respondInteraction(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}

func deferInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	})
}

func editInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, content string) {
	_, _ = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &content,
	})
}
