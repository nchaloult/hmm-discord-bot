package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Starter describes objects which perform necessary procedures before spinning up a Discord bot,
// and then spin up a Discord bot.
//
// I need to come up with a better name for this interface lol
type Starter interface {
	Start()
}

// Bot establishes a new Discord session and is invoked by commands in Discord messages.
type Bot struct {
	dg     *discordgo.Session
	prefix string
	hmm    *HMM
}

// NewBot returns a pointer to a new Bot initialized with the providen token, bot prefix, and hidden
// Markov model to generate content.
func NewBot(token, prefix string, hmm *HMM) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		dg:     dg,
		prefix: prefix,
		hmm:    hmm,
	}, nil
}

// Start opens a websocket connection with Discord and starts listening for events.
func (b *Bot) Start() error {
	b.addHandlers()

	err := b.dg.Open()
	if err != nil {
		return err
	}
	fmt.Println("Bot is up & running. Press Ctrl+C to shut it down.")

	// Block until Ctrl+C is pressed, or the process is interrupted or terminated.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	fmt.Println("Spinning down....")
	b.dg.Close()
	return nil
}

// messageCreateHandler is called every time a new message is posted in a a channel that the bot has
// access to.
func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages posted by the bot.
	// Just to save CPU cycles, even though they're cheap ;)
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong!")
	}
}

// addHandlers registers all of this bot's handler functions with the bot's Discord session.
func (b *Bot) addHandlers() {
	b.dg.AddHandler(messageCreateHandler)
}
