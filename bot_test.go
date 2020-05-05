package main

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
)

// Globals set by postDiscordMessageMock()
var (
	wasMessagePosted bool
	postedMsg        string
)

// TestMessageCreateHandler makes sure that the bot responds appropriately, either with the right
// warning message or with a generated speech that has the properties that were asked for.
//
// This test is long; it might be difficult for someone else other than the original author to read
// and figure out what's going on. I need to think about how to break it up into separate test
// functions without them all being really repetitive.
func TestMessageCreateHandler(t *testing.T) {
	corpus := "the quick brown fox jumps over the lazy dog\n"
	maxRetries := 5
	hmm, _ := NewHMM(corpus, maxRetries)

	botName := "foo"
	botPrefix := "!"
	botToken := "bar"
	bot, _ := NewBot(botName, botPrefix, botToken, hmm)
	bot.postFN = postDiscordMessageMock

	// Setting up test case for message posted by bot.
	botID := "botID"
	s := &discordgo.Session{
		State: &discordgo.State{
			Ready: discordgo.Ready{
				User: &discordgo.User{
					ID: botID,
				},
			},
		},
	}
	m := &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: botID,
			},
		},
	}
	// Make sure that nothing happens, since messages posted by the bot should be ignored.
	bot.MessageCreateHandler(s, m)
	if wasMessagePosted {
		t.Error("A message was posted in response to the bot. MessageCreateHandler should return" +
			" once it realizes that the bot posted the most recent message.")
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case for message that doesn't invoke the bot.
	normalUserID := "normalUserID"
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content: "foo",
		},
	}
	// Make sure nothing happens, since messages that don't invoke the bot should be ignored.
	bot.MessageCreateHandler(s, m)
	if wasMessagePosted {
		t.Error("A message was posted in response to a message that didn't invoke the bot." +
			" Messages that don't invoke the bot should be ignored.")
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case for message that mentions (@s) another user.
	botInvocationString := bot.prefix + bot.name
	mentionedUser := &discordgo.User{}
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content:  botInvocationString + "foo",
			Mentions: []*discordgo.User{mentionedUser},
		},
	}
	// Make sure nothing happens, since messages that mention other users should be ignored.
	bot.MessageCreateHandler(s, m)
	want := "@'ing people isn't supported yet :("
	if !wasMessagePosted {
		t.Error("A message was posted in response to a message that mentioned (@d) another user." +
			" The bot's supposed to post a message notifying the user that they can't do this," +
			" but no such message was posted.")
	} else if postedMsg != want {
		t.Errorf("Unexpected response for message which mentioned a user.\ngot: %q\nwant:%q\n",
			postedMsg, want)
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case where we ask for a message with a specific word count.
	numWordsWant := 42
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content:  botInvocationString + strconv.Itoa(numWordsWant),
			Mentions: []*discordgo.User{},
		},
	}
	// Make sure that the response contains numWordsWant words.
	bot.MessageCreateHandler(s, m)
	if !wasMessagePosted {
		t.Errorf("No message was posted after asking for a speech with %d words.", numWordsWant)
	} else {
		got := len(strings.Fields(postedMsg))
		if got != numWordsWant {
			t.Errorf("Unexpected message response length. got: %d, want:%d\n", got, numWordsWant)
		}
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case where we ask for a message with 0 words.
	numWordsWant = 0
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content:  botInvocationString + strconv.Itoa(numWordsWant),
			Mentions: []*discordgo.User{},
		},
	}
	// Make sure that the response contains numWordsWant words.
	bot.MessageCreateHandler(s, m)
	want = "Can't post an empty message"
	if !wasMessagePosted {
		t.Error("No warning message was posted after asking for a speech with 0 words.")
	} else if postedMsg != want {
		t.Errorf("Unexpected response for message which asked for speech with 0 words."+
			" got: %q, want:%q\n", postedMsg, want)
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case where we ask for a message that begins with a word.
	firstWordWant := "foo"
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content:  botInvocationString + firstWordWant,
			Mentions: []*discordgo.User{},
		},
	}
	// Make sure that the response starts with firstWordWant.
	bot.MessageCreateHandler(s, m)
	if !wasMessagePosted {
		t.Errorf("No message was posted after asking for a speech that begins with %q.",
			firstWordWant)
	} else {
		got := strings.Fields(postedMsg)[0]
		if got != firstWordWant {
			t.Errorf("Message contained unexpected first word. got: %q, want:%q\n",
				got, firstWordWant)
		}
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case where we ask for a message that begins with a word and contains a
	// certain amount of words.
	firstWordWant = "foo"
	numWordsWant = 42
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content:  botInvocationString + firstWordWant + " " + strconv.Itoa(numWordsWant),
			Mentions: []*discordgo.User{},
		},
	}
	// Make sure that the response meets the requested criteria.
	bot.MessageCreateHandler(s, m)
	if !wasMessagePosted {
		t.Errorf("No message was posted after asking for a speech that begins with %q and has %d"+
			" words.", firstWordWant, numWordsWant)
	} else {
		postedMsgWords := strings.Fields(postedMsg)
		firstWordGot := postedMsgWords[0]
		numWordsGot := len(postedMsgWords)
		if firstWordGot != firstWordWant || numWordsGot != numWordsWant {
			t.Errorf("Message did not meet requested criteria. got: %q and %d, want: %q and %d\n",
				firstWordGot, numWordsGot, firstWordWant, numWordsWant)
		}
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case where we ask for a message that begins with a word and contains 0 words.
	firstWordWant = "foo"
	numWordsWant = 0
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content:  botInvocationString + firstWordWant + " " + strconv.Itoa(numWordsWant),
			Mentions: []*discordgo.User{},
		},
	}
	// Make sure that the response meets the requested criteria.
	bot.MessageCreateHandler(s, m)
	want = "Can't post an empty message"
	if !wasMessagePosted {
		t.Error("No warning message was posted after asking for a speech that begins with a word" +
			" and has 0 words.")
	} else if postedMsg != want {
		t.Errorf("Unexpected response for message which asked for speech that begins with a word"+
			" and has 0 words. got: %q, want: %q\n", postedMsg, want)
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case where we make a request with too many arguments.
	firstWordWant = "foo"
	numWordsWant = 42
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content: botInvocationString + firstWordWant + " " + strconv.Itoa(numWordsWant) + " " +
				"bar baz",
			Mentions: []*discordgo.User{},
		},
	}
	// Make sure that the response meets the requested criteria, ignoring the additional arguments.
	bot.MessageCreateHandler(s, m)
	if !wasMessagePosted {
		t.Errorf("No message was posted after asking for a speech that begins with %q and has %d"+
			" words (but with too many arguments).", firstWordWant, numWordsWant)
	} else {
		postedMsgWords := strings.Fields(postedMsg)
		firstWordGot := postedMsgWords[0]
		numWordsGot := len(postedMsgWords)
		if firstWordGot != firstWordWant || numWordsGot != numWordsWant {
			t.Errorf("Message did not meet requested criteria. got: %q and %d, want: %q and %d\n",
				firstWordGot, numWordsGot, firstWordWant, numWordsWant)
		}
	}
	wasMessagePosted = false
	postedMsg = ""

	// Setting up test case where we make an invalid request with two arguments (firstWord and
	// numWords are flipped).
	m = &discordgo.MessageCreate{
		Message: &discordgo.Message{
			Author: &discordgo.User{
				ID: normalUserID,
			},
			Content:  botInvocationString + strconv.Itoa(numWordsWant) + " " + firstWordWant,
			Mentions: []*discordgo.User{},
		},
	}
	// Make sure that the response contains the appropriate warning.
	bot.MessageCreateHandler(s, m)
	want = fmt.Sprintf("%q is not a number. Example usage: `%s <firstWord> <numWords>`",
		firstWordWant, botInvocationString)
	if !wasMessagePosted {
		t.Errorf("No warning message was posted after asking for a speech that begins with %q and"+
			" has %d words, but with the arguments flipped.", firstWordWant, numWordsWant)
	} else if postedMsg != want {
		t.Errorf("Unexpected response for message with invalid arguments.\ngot: %q\nwant:%q\n",
			postedMsg, want)
	}
	wasMessagePosted = false
	postedMsg = ""
}

// postDiscordMessageMock is a function of the custom type: MsgPoster. It sets the global:
// "wasMessagePosted" and "msg" vars. A test is then responsible for checking the values of those
// global vars and resetting them.
func postDiscordMessageMock(session *discordgo.Session, channelID, msg string) {
	wasMessagePosted = true
	postedMsg = msg
}
