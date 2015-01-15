package sshtunnel

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"net/http"
)

type (
	Config struct {
		SshHost  string
		SshPort  int
		TpHost   string
		TpPort   int
		Username string
		Password string
	}

	Tunnel struct {
		client *ssh.Client
		tpHost string
		tpPort int
	}
)

func Connect(cfg *Config) (*Tunnel, error) {
	sshcfg := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Password),
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.SshHost, cfg.SshPort), sshcfg)

	if err != nil {
		return nil, err
	}

	return &Tunnel{client, cfg.TpHost, cfg.TpPort}, nil
}

func (t *Tunnel) Get(req *http.Request) (*http.Response, error) {
	tunnel, err := t.client.Dial("tcp", fmt.Sprintf("%s:%d", t.tpHost, t.tpPort))

	if err != nil {
		return nil, err
	}

	tunnelReader := bufio.NewReader(tunnel)

	req.Write(tunnel)

	response, err := http.ReadResponse(tunnelReader, req)

	if err != nil {
		return nil, err
	}

	return response, nil
}
