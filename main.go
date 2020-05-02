package main

import (
	"bufio"
	"fmt"
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
		log.Fatalf("Failed to read corpus file: %v\n", err)
	}
	defer file.Close()

	s := bufio.NewScanner(file)
	numLines := 10
	for i := 0; i < numLines; i++ {
		s.Scan()
		fmt.Println(s.Text())
	}
}
