package handlers

import (
	"github.com/bwmarrin/discordgo"
)

func BoostHandler(
	session *discordgo.Session,
	event *discordgo.GuildMemberUpdate,
	roleBoost string,
	pendaRole string,
	pendaGoldRole string,
) error {
	guildID := event.GuildID
	memberID := event.User.ID
	member := event.Member

	//member, err := session.GuildMember(guildID, memberID)
	//if err != nil {
	//	return err
	//}

	hasPenda := false
	for _, role := range member.Roles {
		if role == pendaRole {
			hasPenda = true
			break
		}
	}

	if hasPenda {
		hasBoost := false
		for _, role := range member.Roles {
			if role == roleBoost {
				hasBoost = true
				break
			}
		}

		if hasBoost {
			err := session.GuildMemberRoleAdd(guildID, memberID, pendaGoldRole)
			if err != nil {
				return err
			}
		} else {
			err := session.GuildMemberRoleRemove(guildID, memberID, pendaGoldRole)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
