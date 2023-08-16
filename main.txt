package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file", err)
	}

	signer, err := LoadPrivKey()
	if err != nil {
		log.Fatal("failed to laod private key: ", err)
	}
	client := NewSSHClient(signer)
	output, err := client.RunCmd("ls -l")
	if err != nil {
		log.Fatal("failed to run command: ", err)
	}
	fmt.Println(output)
}
