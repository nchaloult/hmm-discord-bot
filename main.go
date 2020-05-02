package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
)

const corporaDirName = "corpora"

func main() {
	botname := os.Getenv("BOT_NAME")
	prefix := os.Getenv("BOT_PREFIX")
	token := os.Getenv("BOT_TOKEN")
	filename := os.Getenv("FILENAME")

	// Read the corpus file and "train" a hidden Markov model.
	file, err := os.Open(path.Join(corporaDirName, filename))
	if err != nil {
		log.Fatalf("Failed to open corpus file: %v\n", err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read corpus file: %v\n", err)
	}
	hmm := NewHMM(string(content), 5)

	// Create a Discord bot and spin it up.
	bot, err := NewBot(botname, prefix, token, hmm)
	if err != nil {
		log.Fatalf("Failed to create new Discord bot: %v\n", err)
	}
	err = bot.Start()
	if err != nil {
		log.Fatalf("Failed to spin up Discord bot: %v\n", err)
	}
}
