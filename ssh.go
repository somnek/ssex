package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	client *ssh.Client
}

func LoadPrivKey() (ssh.Signer, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("faild to get user home directory: ", err)
	}

	key, err := os.ReadFile(home + "/.ssh/id_rsa")
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return signer, nil
}

func NewSSHClient(signer ssh.Signer) *Client {
	user := os.Getenv("USERNAME")
	host := os.Getenv("HOST")

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	client, err := ssh.Dial("tcp", host, config)
	if err != nil {
		log.Fatal("failed at dial:", err)
	}
	return &Client{client: client}
}

func (c *Client) Close() {
	c.client.Close()
}

func (c *Client) RunCmd(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		return "", fmt.Errorf("failed to run command: %v", err)
	}
	return b.String(), nil
}
