package main

import (
	"bufio"
	"code.google.com/p/gcfg"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net/http"
	"os"
)

type (
	Config struct {
		Ssh     sshSection
		Tikplay tikplaySection
	}

	sshSection struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	tikplaySection struct {
		Host string
		Port int
	}
)

var (
	defaultConfig = Config{
		sshSection{Host: "kekkonen.cs.hut.fi", Port: 22},
		tikplaySection{Host: "tikradio.tt.hut.fi", Port: 5000},
	}
)

func handleError(e error, s string, f bool) {
	if e != nil {
		fmt.Printf("%s: %s\n", s, e.Error())
		if f {
			os.Exit(1)
		}
	}
}

func main() {
	cfg := defaultConfig
	err := gcfg.ReadFileInto(&cfg, "config.gcfg")
	handleError(err, "config", true)

	sshcfg := &ssh.ClientConfig{
		User: cfg.Ssh.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(cfg.Ssh.Password),
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Ssh.Host, cfg.Ssh.Port), sshcfg)

	handleError(err, "ssh client", true)

	tunnel, err := client.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Tikplay.Host, cfg.Tikplay.Port))

	handleError(err, "tunnel", true)

	tunnelReader := bufio.NewReader(tunnel)

	request, err := http.NewRequest("GET", "/srv/v1.0/song", nil)
	handleError(err, "Request", true)

	request.Write(tunnel)
	response, err := http.ReadResponse(tunnelReader, request)
	handleError(err, "Response", true)
	if response != nil {
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		handleError(err, "Body", true)
		fmt.Printf("%s\n", body)
	}

}
