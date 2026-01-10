package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"pendago/handlers"
	"pendago/modules"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	BotToken  = flag.String("t", "", "Bot token")
	BotConfig = flag.String("c", "", "JSON config file")
)

func init() { flag.Parse() }

func main() {
	log.Println("Given file : ", *BotConfig)
	err := modules.ReadConfig(*BotConfig)
	if err != nil {
		log.Println("Fail while reading config file ", err)
	}

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
	discordBot.AddHandler(interactionCreate)

	_, err = discordBot.ApplicationCommandCreate(
		discordBot.State.User.ID,
		"",
		&discordgo.ApplicationCommand{
			Name:        "minecraft",
			Description: "Commandes Minecraft.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionSubCommand,
					Name:        "whitelist_add",
					Description: "Ajoute un joueur à la whitelist Minecraft.",
					Options: []*discordgo.ApplicationCommandOption{
						{
							Type:        discordgo.ApplicationCommandOptionString,
							Name:        "pseudo",
							Description: "Pseudo Minecraft",
							Required:    true,
						},
					},
				},
			},
		},
	)
	if err != nil {
		log.Println("Error creating command: ", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit or kill processus.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	discordBot.Close()
}

func guildMemberBoost(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	err := handlers.BoostHandler(
		s,
		event,
		modules.RoleBoost,
		modules.PendaRole,
		modules.PendaGoldRole,
	)
	if err != nil {
		log.Fatal("Error during role attribution : ", err)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	handlers.MessageHandler(
		s,
		m,
		modules.WowLogChanId,
		modules.WowChanId,
		modules.BotChanId,
		modules.TriggerCommand,
		modules.RoleId,
		modules.SWCRoleId,
		modules.SWCommand,
	)
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	handlers.InteractionHandler(
		s,
		i,
		modules.BotChanId,
		modules.PendaRole,
		modules.RCONHost,
		modules.RCONPort,
		modules.RCONPassword,
	)
}
