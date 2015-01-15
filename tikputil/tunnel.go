package tikputil

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"net/http"
	"os"
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

	ConnWrap struct {
		conn net.Conn
	}
)

func (cw ConnWrap) Read(p []byte) (n int, err error) {
	defer fmt.Println("R:", p)
	return cw.conn.Read(p)
}
func (cw ConnWrap) Write(p []byte) (n int, err error) {
	defer fmt.Println("W:", p)
	return cw.conn.Write(p)
}

func Connect(cfg *Config) (*Tunnel, error) {
	auths := []ssh.AuthMethod{ssh.Password(cfg.Password)}

	if auth := getSignersFromAgent(); auth != nil {
		auths = append(auths, auth)
	}

	sshcfg := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: auths,
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.SshHost, cfg.SshPort), sshcfg)

	if err != nil {
		return nil, err
	}

	return &Tunnel{client, cfg.TpHost, cfg.TpPort}, nil
}

func getSignersFromAgent() ssh.AuthMethod {
	sock, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		return nil
	}
	agent := agent.NewClient(sock)
	signers, err := agent.Signers()
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(signers...)
}

func (t *Tunnel) Execute(req *http.Request) (*http.Response, error) {
	tunnel, err := t.client.Dial("tcp", fmt.Sprintf("%s:%d", t.tpHost, t.tpPort))

	if err != nil {
		return nil, err
	}

	//defer tunnel.Close()

	mock := ConnWrap{tunnel}

	tunnelReader := bufio.NewReader(mock)

	req.Write(mock)

	response, err := http.ReadResponse(tunnelReader, req)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func (t *Tunnel) Close() {
	t.client.Conn.Close()
}
