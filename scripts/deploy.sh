#!/bin/bash

heroku container:push worker -a hmm-discord-bot
heroku container:release worker -a hmm-discord-bot
