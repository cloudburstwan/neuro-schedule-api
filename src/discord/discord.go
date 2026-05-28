package discord

import (
	"errors"
	"fmt"
	"neuro-schedule-api/src/hub"
	"neuro-schedule-api/src/parser"

	"github.com/bwmarrin/discordgo"
)

var discord *discordgo.Session
var listenChannelId string
var hubPointer *hub.Hub

// StartDiscordListener attempts to start the Discord component of the schedule parser
// which listens for schedule updates in the schedule channel.
func StartDiscordListener(token string, channelId string, hub *hub.Hub) error {
	listenChannelId = channelId
	hubPointer = hub

	var err error
	discord, err = discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Encountered an error while creating Discord session: ", err)
		return errors.New("error creating Discord session")
	}

	discord.AddHandler(onMessageCreate)
	discord.AddHandler(onMessageUpdate)

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages | discordgo.IntentsMessageContent

	err = discord.Open()
	if err != nil {
		fmt.Println("Encountered an error while opening Discord session: ", err)
		return errors.New("error opening Discord session")
	}
	return nil
}

// StopDiscordListener closes the Discord bot session.
func StopDiscordListener() error {
	err := discord.Close()
	if err != nil {
		fmt.Println("Encountered an error while closing Discord session: ", err)
		return errors.New("error closing Discord session")
	}
	return nil
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID != listenChannelId {
		return
	}
	fmt.Println("Received a messageCreate in schedule channel. Running schedule parser...")

	parsedSchedule, err := parser.Parse(m.Content, hubPointer.GetSchedule())
	if err != nil {
		fmt.Println("Encountered an error while parsing schedule from Discord messageCreate: ", err)
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		if err != nil {
			return
		}
		return
	}

	fmt.Println("Schedule parsed, updating saved schedule information...")

	err = hubPointer.Broadcast(parsedSchedule)
	if err != nil {
		fmt.Println("Encountered an error while updating schedule due to Discord messageCreate: ", err)
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		if err != nil {
			return
		}
		return
	}

	fmt.Println("Schedule parsed and saved.")

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		return
	}
}

func onMessageUpdate(s *discordgo.Session, m *discordgo.MessageUpdate) {
	if m.ChannelID != listenChannelId {
		return
	}
	fmt.Println("Received a messageUpdate in schedule channel. Running schedule parser...")

	parsedSchedule, err := parser.Parse(m.Content, hubPointer.GetSchedule())
	if err != nil {
		fmt.Println("Encountered an error while parsing schedule from Discord messageUpdate: ", err)
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		if err != nil {
			return
		}
		return
	}

	fmt.Println("Schedule parsed, updating saved schedule information...")

	err = hubPointer.Broadcast(parsedSchedule)
	if err != nil {
		fmt.Println("Encountered an error while updating schedule due to Discord messageUpdate: ", err)
		err := s.MessageReactionAdd(m.ChannelID, m.ID, "🚫")
		if err != nil {
			return
		}
		return
	}

	fmt.Println("Schedule parsed and saved.")

	err = s.MessageReactionAdd(m.ChannelID, m.ID, "✅")
	if err != nil {
		return
	}
}
