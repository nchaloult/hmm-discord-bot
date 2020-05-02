# hmm-discord-bot

A Discord bot that generates messages of the same vocabulary and sentence structure as a provided
corpus.

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
