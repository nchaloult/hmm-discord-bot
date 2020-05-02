package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	_ "github.com/joho/godotenv/autoload"
)

const corporaDirName = "corpora"

func main() {
	filename := os.Getenv("FILENAME")

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

	fmt.Println(hmm.GenerateSpeech())
}
