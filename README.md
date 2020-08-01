# hmm-discord-bot

[![Build Status](https://travis-ci.org/nchaloult/hmm-discord-bot.svg?branch=master)](https://travis-ci.org/nchaloult/hmm-discord-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/nchaloult/hmm-discord-bot)](https://goreportcard.com/report/github.com/nchaloult/hmm-discord-bot)

A Discord bot that generates messages of the same vocabulary and sentence structure as a provided
corpus.

A [hidden Markov model](https://en.wikipedia.org/wiki/Hidden_Markov_model) is used to generate
messages. For speech generation, HMMs are notoriously mediocre, but I think that's part of the fun!
:) This bot can pump out some hilarious garbage sometimes.

<img width="818" alt="Screen shot that showcases the Discord bot in action" src="https://user-images.githubusercontent.com/31291920/89092797-efc38980-d382-11ea-98cd-5e65949a9671.png">

## Supported Commands

When invoking the bot with the configured name and prefix, it supports the following argument patterns:

- No arguments: generates a message of unknown length until the HMM decides that the message should end
    - Ex: `!botname`
- `<num-words>`: generates a message with the provided number of words
    - Ex: `!botname 40`
- `<beginning-word>`: generates a message that starts with the provided word
    - Ex: `!botname america`
    - The provided `<beginning-word>` does NOT need to be in the corpus file that the HMM is trained on, although the results you get are often better if it is
- `<beginning-word> <num-words>`: generates a message with the provided number of words AND that starts with the provided word
    - Ex: `!botname america 40`

## Configuration

The bot may be configured with a name and a prefix. These two things are what users type in Discord messages to invoke the bot. For instance, in the screenshot above, the bot's configured name is "obama", and its prefix is "!".

The bot also needs to be configured with the name of a corpus file to train an HMM on, as well as a Discord API token. That corpus file needs to live in the `/corpora` directory. Instructions for provisioning an API token for a Discord bot can be found [here](https://discordpy.readthedocs.io/en/latest/discord.html).

All four of those configurable items are kept in environment variables. If you want to deploy an instance of this bot and bring it into a Discord server that you're a part of, you'll need to set those environment variables in whatever deployment environment you end up working with. See the [`.env.sample`](.env.sample) file for which environment variables you'll need to set.

## Development Setup

1. Clone this repo
1. Create a `.env` file in the root dir of this project
    * Either `$ mv .env.sample .env`
        * If you do this, then you'll need to replace the `BOT_TOKEN` var with your own bot's token
    * Or `$ touch .env` and fill that file up yourself
1. Place your corpus file(s) in the `/corpora` dir
    * See the README in `/corpora` for more info about corpus files
    * tl;dr: an example corpus file is provided if you just wanna look at that or use it
1. Spin up the program
    * `$ go run *.go`
    * Or `$ go build && ./hmm-discord-bot` if you're hungry for speed

## On Deploying to Production

This application is containerized, so you should be able to deploy this bot wherever you can host and deploy containers 🤞
