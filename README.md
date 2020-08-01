# hmm-discord-bot

[![Build Status](https://travis-ci.org/nchaloult/hmm-discord-bot.svg?branch=master)](https://travis-ci.org/nchaloult/hmm-discord-bot)
[![Go Report Card](https://goreportcard.com/badge/github.com/nchaloult/hmm-discord-bot)](https://goreportcard.com/report/github.com/nchaloult/hmm-discord-bot)

A Discord bot that generates messages of the same vocabulary and sentence structure as a provided
corpus.

A [hidden Markov model](https://en.wikipedia.org/wiki/Hidden_Markov_model) is used to generate
messages. For speech generation, HMMs are notoriously mediocre, but I think that's part of the fun!
:) This bot can pump out some hilarious garbage sometimes.

<img width="818" alt="Screen shot that showcases the Discord bot in action" src="https://user-images.githubusercontent.com/31291920/89092797-efc38980-d382-11ea-98cd-5e65949a9671.png">

## Initial Setup

1. Create a `.env` file in the root dir of this project
    * Either `$ mv .env.sample .env`
        * If you do this, then you'll need to replace the `BOT_TOKEN` var with your own bot's token
    * Or `$ touch .env` and fill it up yourself
1. Place your corpus file(s) in the `/corpora` dir
    * See the README in `/corpora` for more info about corpus files
    * tl;dr: an example corpus file is provided if you just wanna look at that or use it
1. Spin up the program
    * `$ go run *.go`
    * Or `$ go build && ./hmm-discord-bot` if you're hungry for speed
