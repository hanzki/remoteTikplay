package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/hanzki/remoteTikplay/sshtunnel"
	"github.com/hanzki/remoteTikplay/tikplay"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type (
	Config struct {
		Tunnel sshtunnel.Config
	}
)

var (
	defaultConfig = Config{
		sshtunnel.Config{
			SshHost: "kekkonen.cs.hut.fi",
			SshPort: 22,
			TpHost:  "tikradio.tt.hut.fi",
			TpPort:  5000,
		},
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

	tunnel, err := sshtunnel.Connect(&cfg.Tunnel)

	handleError(err, "ssh client", true)

	defer tunnel.Close()

	var request *http.Request

	switch os.Args[1] {
	case "np":
		request, err = tikplay.NowPlaying()
	case "skip":
		request, err = tikplay.Skip()
	case "clear":
		request, err = tikplay.Clear()
	case "list":
		n, err := strconv.Atoi(os.Args[2])
		if err == nil {
			request, err = tikplay.Playlist(uint(n))
		}
	case "play":
		request, err = tikplay.Play(
			&tikplay.PlayJSON{
				fmt.Sprintf("%s@%s", cfg.Tunnel.Username, cfg.Tunnel.SshHost),
				os.Args[2],
			})
	}

	handleError(err, "Request", true)

	response, err := tunnel.Execute(request)

	handleError(err, "Response", true)

	if response != nil {
		body, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		handleError(err, "Body", true)
		fmt.Printf("%s\n", body)
	}
}
