package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"strings"
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
	dg            *discordgo.Session
	name          string
	prefix        string
	hmm           *HMM
	contentRegexp *regexp.Regexp
}

// NewBot returns a pointer to a new Bot initialized with the providen token, bot prefix, and hidden
// Markov model to generate content.
func NewBot(name, prefix, token string, hmm *HMM) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// This step is kinda expensive. Instead of doing this in messageCreateHandler() every time we
	// handle a bot invocation, we do this once when the bot is created.
	reg, err := regexp.Compile("[^a-zA-Z0-9 ]+") // Filtering for alphanumeric values only
	if err != nil {
		return nil, err
	}

	return &Bot{
		dg:            dg,
		name:          name,
		prefix:        prefix,
		hmm:           hmm,
		contentRegexp: reg,
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

	fmt.Println("\nSpinning down....")
	b.dg.Close()
	return nil
}

// messageCreateHandler is called every time a new message is posted in a a channel that the bot has
// access to.
func (b *Bot) messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages posted by the bot.
	// Just to save CPU cycles, even though they're cheap ;)
	if m.Author.ID == s.State.User.ID {
		return
	}

	prefixAndName := b.prefix + b.name
	if !strings.HasPrefix(m.Content, prefixAndName) {
		return
	}

	if len(m.Mentions) > 0 {
		s.ChannelMessageSend(m.ChannelID, "@'ing people isn't supported yet :(")
		return
	}

	// Input cleanup and sanitization
	content := strings.TrimPrefix(m.Content, prefixAndName)
	content = strings.TrimSpace(content)
	content = strings.ToLower(content)
	content = b.contentRegexp.ReplaceAllString(content, "")
	// Just echo content for now
	s.ChannelMessageSend(m.ChannelID, content)
}

// addHandlers registers all of this bot's handler functions with the bot's Discord session.
func (b *Bot) addHandlers() {
	b.dg.AddHandler(b.messageCreateHandler)
}
