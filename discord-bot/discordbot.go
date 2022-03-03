package discordbot

import (
	"les-randoms/discord-bot/logic"
	"les-randoms/utils"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var appEnd *chan bool // If a value is sent to this, the complete app stop

var Bot *logic.DiscordBot

func Start(applicationEnd *chan bool) {
	appEnd = applicationEnd

	var err error
	err = Bot.Start()
	if err != nil {
		utils.LogError(err.Error())
		return
	}

	Bot.SetDefaultCommand(CommandUnknown)
	Bot.AddCommand("SHUTDOWN", CommandShutdown)
	Bot.AddCommand("PING", CommandPing)
	Bot.AddCommand("PLAY", CommandPlay)
	Bot.AddCommand("KANNA", func(bot *logic.DiscordBot, m *discordgo.MessageCreate) error {
		return nil
	})
	Bot.ActivateMessageHandler(func(bot *logic.DiscordBot, m *discordgo.MessageCreate) string {
		return strings.Split(strings.ToUpper(m.Content), " ")[0]
	})

	utils.LogSuccess("Discord bot successfully started")
}

func Close() {
	utils.LogNotNilError(Bot.DiscordSession.Close())
	utils.LogSuccess("Discord bot session successfully closed")
}
