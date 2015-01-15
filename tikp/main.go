package main

import (
	"code.google.com/p/gcfg"
	"fmt"
	"github.com/hanzki/remoteTikplay/tikputil"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type (
	Config struct {
		Tunnel tikputil.Config
	}
)

const usage = `usage: file <command> <parameter>
available commands:
np         = now playing
list <n>   = lists n songs from queue (n defaults to 10)
play <url> = plays song from url
skip       = skips the current song
clear      = clears the whole queue`

var (
	defaultConfig = Config{
		tikputil.Config{
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
	if len(os.Args) < 2 {
		fmt.Println(usage)
		os.Exit(0)
	}

	cfg := defaultConfig
	err := gcfg.ReadFileInto(&cfg, "config.gcfg")
	handleError(err, "config", true)

	tunnel, err := tikputil.Connect(&cfg.Tunnel)

	handleError(err, "ssh client", true)

	defer tunnel.Close()

	var request *http.Request

	switch os.Args[1] {
	case "np":
		request, err = tikputil.NowPlaying()
	case "skip":
		request, err = tikputil.Skip()
	case "clear":
		request, err = tikputil.Clear()
	case "list":
		var (
			n   int   = 10
			err error = nil
		)
		if len(os.Args) >= 3 {
			n, err = strconv.Atoi(os.Args[2])
		}
		if err == nil {
			request, err = tikputil.Playlist(uint(n))
		}
	case "play":
		if len(os.Args) >= 3 {
			request, err = tikputil.Play(
				&tikputil.PlayJSON{
					fmt.Sprintf("%s@%s", cfg.Tunnel.Username, cfg.Tunnel.SshHost),
					os.Args[2],
				})
		} else {
			fmt.Println("Missing play url")
			os.Exit(1)
		}
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
