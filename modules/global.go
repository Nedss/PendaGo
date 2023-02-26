package modules

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func createPollEmbedMessage(
	question string,
	answer string,
	answerNumber string,
) *discordgo.MessageEmbed {
	pollEmbed := NewEmbed().
		SetTitle(question).
		SetDescription(answer).
		AddField("Nombre de réponses", answerNumber).
		SetColor(0x66E36E).MessageEmbed
	return pollEmbed
}

func splitPollMessage(message string) ([]string, error) {
	splitMessage := strings.Split(message, "\"")

	// Remove empty
	var formatedSplitMessage []string
	for _, value := range splitMessage {
		// Only accept non traillig white spaced string
		if len(strings.TrimSpace(value)) > 0 {
			formatedSplitMessage = append(formatedSplitMessage, value)
		}
	}

	if len(formatedSplitMessage) < 4 {
		return []string{}, fmt.Errorf("Message is empty, can't be splitted")
	}

	return formatedSplitMessage, nil
}

func generatePollEmbedMessage(m *discordgo.MessageCreate) (*discordgo.MessageEmbed, error) {
	emoji := [9]string{
		"one",
		"two",
		"three",
		"four",
		"five",
		"six",
		"seven",
		"eight",
		"nine",
	}

	finalMessage := ""

	splitAnswer, err := splitPollMessage(m.Content)
	if err != nil {
		return nil, err
	}

	question := splitAnswer[1]
	splitAnswer = splitAnswer[2:]

	answerNumber := len(splitAnswer)

	if answerNumber == 0 {
		return nil, fmt.Errorf("No answer is given in the poll.")
	}

	for index, element := range splitAnswer {
		finalMessage = finalMessage + ":" + emoji[index] + ":" + " " + element + "\n"
	}

	embedPoll := createPollEmbedMessage(question, finalMessage, strconv.Itoa(answerNumber))
	return embedPoll, nil
}

func GeneratePollEmbedReaction(s *discordgo.Session, m *discordgo.MessageCreate) error {
	umoji := [9]string{
		"\u0031\u20E3",
		"\u0032\u20E3",
		"\u0033\u20E3",
		"\u0034\u20E3",
		"\u0035\u20E3",
		"\u0036\u20E3",
		"\u0037\u20E3",
		"\u0038\u20E3",
		"\u0039\u20E3",
	}

	channelId := m.ChannelID

	me, err := generatePollEmbedMessage(m)
	if err != nil {
		return err
	}

	sentMessage, err := s.ChannelMessageSendEmbed(channelId, me)
	if err != nil {
		return err
	}

	sentMessageId := sentMessage.ID
	answerLength, err := strconv.Atoi(me.Fields[0].Value)
	if err != nil {
		return err
	}

	for i := 0; i < answerLength; i++ {
		s.MessageReactionAdd(channelId, sentMessageId, umoji[i])
	}

	return nil
}

func BoostHandler(
	s *discordgo.Session,
	e *discordgo.GuildMemberUpdate,
	boostRoleID string,
	pendaRoleID string,
	pendaGoldRole string,
) error {
	if e.Member != nil && e.Member.Roles != nil {
		var isPenda bool = false
		var isBooster bool = false
		for _, roleID := range e.Member.Roles {
			if roleID == boostRoleID {
				isBooster = true
			}
			if roleID == pendaRoleID {
				isPenda = true
			}
		}
		if isPenda == true && isBooster == true {
			err := s.GuildMemberRoleAdd(e.GuildID, e.User.ID, pendaGoldRole)
			if err != nil {
				log.Fatal("Error during role attribution to member", err)
				return err
			}
			return nil
		}
	} else {
		return errors.New("Wrong member and Roles")
	}
	return nil
}
