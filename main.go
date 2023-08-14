package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/ssh"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	user := os.Getenv("USERNAME")
	host := os.Getenv("HOST")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Failed to get user home directory: ", err)
	}

	key, err := os.ReadFile(home + "/.ssh/id_rsa")
	if err != nil {
		log.Fatal("Failed to load private key: ", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("Failed to parse private key: ", err)
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatal("Failed at dial:", err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to created session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("ls -l"); err != nil {
		log.Fatal("Failed to run: ", err.Error())
	}
	fmt.Println(b.String())
}
