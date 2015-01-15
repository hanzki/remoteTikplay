package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/hanzki/remoteTikplay/sshtunnel"
	"io/ioutil"
	"net/http"
	"os"
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

	for i := 0; i < 10; i++ {
		request, err := http.NewRequest("GET", "/srv/v1.0/song", nil)

		handleError(err, "Request", true)

		response, err := tunnel.Get(request)

		handleError(err, "Response", true)

		if response != nil {
			body, err := ioutil.ReadAll(response.Body)
			response.Body.Close()
			handleError(err, "Body", true)
			fmt.Printf("%s\n", body)
		}
	}

}
