package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// msgTooLongErr is a common error that may occur when attempting to post a Discord message. This
// error is checked for in postDiscordMessage().
const msgTooLongErr = "HTTP 400 Bad Request, {\"content\": [\"Must be 2000 or fewer in length.\"]}"

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
	postFN        MsgPoster
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
		postFN:        postDiscordMessage,
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
	log.Println("Bot is up & running. Press Ctrl+C to shut it down.")

	// Block until Ctrl+C is pressed, or the process is interrupted or terminated.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println("Spinning down....")
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
	// Look for bot prefix at the beginning of the message.
	prefixAndName := b.prefix + b.name
	if !strings.HasPrefix(m.Content, prefixAndName) {
		return
	}
	// If anyone was mentioned in the message, don't mess with it.
	if len(m.Mentions) > 0 {
		b.postFN(s, m.ChannelID, "@'ing people isn't supported yet :(")
		return
	}

	// Clean up and sanitize input.
	content := strings.TrimPrefix(m.Content, prefixAndName)
	content = strings.TrimSpace(content)
	content = strings.ToLower(content)
	content = b.contentRegexp.ReplaceAllString(content, "")

	arguments := strings.Split(content, " ")
	numArgs := len(arguments)
	// If content is an empty string, then arguments will look like: [""]. Nuke that empty string.
	if numArgs == 1 && arguments[0] == "" {
		arguments = nil
		numArgs = 0
	}

	// Handle response based on how many arguments were provided in the bot invocation.
	if numArgs == 0 {
		msg := b.hmm.GenerateSpeech()
		b.postFN(s, m.ChannelID, msg)
		return
	}
	if numArgs == 1 {
		arg := arguments[0]
		// Determine if a first word was provided or if a number of words was provided.
		numWords, err := strconv.Atoi(arg)
		if err != nil {
			// Something went wrong trying to convert the first argument to an int. That means the
			// first argument is a word that the generated text should start with.
			msg := b.hmm.GenerateSpeechBeginningWithWord(arg)
			b.postFN(s, m.ChannelID, msg)
			return
		}
		// The string to int conversion was successful. Assume that the number passed in is the
		// number of words that the generated text should have.
		msg := b.hmm.GenerateSpeechWithNumWords(numWords)
		b.postFN(s, m.ChannelID, msg)
		return
	}
	// len(arguments) is at least 2. If there were more than 2 arguments provided, ignore all of
	// them except for the first two.
	firstWord := arguments[0]
	numWords, err := strconv.Atoi(arguments[1])
	if err != nil {
		// Second argument was not a number. Respond with usage instructions.
		msg := fmt.Sprintf("\"%s\" is not a number. Example usage: `%s"+
			" <firstWord> <numWords>`", arguments[1], prefixAndName)
		b.postFN(s, m.ChannelID, msg)
		return
	}
	msg := b.hmm.GenerateSpeechBeginningWithWordAndWithNumWords(firstWord, numWords)
	b.postFN(s, m.ChannelID, msg)
}

// MsgPoster describes functions that send messages to specified Discord channels. This type exists
// mainly so that postDiscordMessage() can be mocked in tests.
type MsgPoster func(*discordgo.Session, string, string)

// postDiscordMessage posts a message in the provided Discord channel as the bot. If anything goes
// wrong, this function is responsible for handling that problem, most likely by just posting a
// message in Discord about what happened.
//
// postDiscordMessage is of the custom type: MsgPoster
func postDiscordMessage(session *discordgo.Session, channelID, msg string) {
	_, err := session.ChannelMessageSend(channelID, msg)
	if err != nil {
		// See if the error is a common one that we recognize.
		if err.Error() == msgTooLongErr {
			errMsg := "The generated message was too long. Discord doesn't let messages that are" +
				" longer than 2000 characters go through."
			session.ChannelMessageSend(channelID, errMsg)
			return
		}
		// We didn't recognize the error at this point. Post a general response.
		errMsg := fmt.Sprintf("Something went wrong: %v", err)
		session.ChannelMessageSend(channelID, errMsg)
	}
}

// addHandlers registers all of this bot's handler functions with the bot's Discord session.
func (b *Bot) addHandlers() {
	b.dg.AddHandler(b.messageCreateHandler)
}
