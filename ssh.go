package main

import (
	"context"
	"fmt"
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
func NewSSHClient(ctx context.Context, signer ssh.Signer, p Profile, resultChan chan<- sshResult) {
	user := p.User
	host := p.Host
	port := p.Port

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.HostKeyCallback(
			func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil },
		),
	}

	// default to port 22
	if port == "" {
		port = "22"
	}
	addr := fmt.Sprintf("%s:%s", host, port)

	select {
	case <-ctx.Done():
		resultChan <- sshResult{err: ctx.Err()}
		return
	default:
		client, err := ssh.Dial("tcp", addr, config)
		if err != nil {
			resultChan <- sshResult{err: fmt.Errorf("failed to dial: %v", err)}
			return
		}
		resultChan <- sshResult{client: &Client{client: client}}
	}
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
		fmt.Println(err)
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
