package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/kevinburke/ssh_config"
	"golang.org/x/crypto/ssh"
)

type Client struct {
	client *ssh.Client
}

// LoadPrivKey loads private key from ~/.ssh/id_rsa
// and returns ssh.Signer interface
func LoadPrivKey() (ssh.Signer, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("faild to get user home directory: %v", err)
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

// NewSSHClient creates ssh client with ssh.Signer interface
// and returns Client struct
func NewSSHClient(signer ssh.Signer, user, addr string) *Client {
	// user = os.Getenv("USERNAME")
	// host = os.Getenv("HOST")
	log.Println("adress: ", addr)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(
			func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
		),
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		fmt.Println("failed at dial:", err)
		os.Exit(1)

	}
	return &Client{client: client}
}

// Close closes ssh client
func (c *Client) Close() {
	c.client.Close()
}

// RunCmd runs command on remote host
// and returns stdout and stderr
func (c *Client) RunCmd(command string) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	data, err := session.CombinedOutput(command)
	return string(data), err
}

func ParseSSHConfig() []*ssh_config.Host {
	f, _ := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	cfg, _ := ssh_config.Decode(f)
	return cfg.Hosts
}
