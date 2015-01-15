package tikputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseUrl  = "/srv/v1.0"
	songUrl  = baseUrl + "/song"
	queueUrl = baseUrl + "/queue"
	fileUrl  = baseUrl + "/file"
	taskUrl  = baseUrl + "/task"
)

type playJSON struct {
	User string `json:"user"`
	Url  string `json:"url"`
}

type Tikplay struct {
	Tunnel *Tunnel
	Whoami string
}

func NewTikplay(cfg *Config) (*Tikplay, error) {
	tunnel, err := Connect(cfg)
	if err != nil {
		return nil, err
	}
	return &Tikplay{
		tunnel,
		fmt.Sprintf("%s@%s", cfg.Username, cfg.SshHost),
	}, nil
}

func (tp *Tikplay) NowPlaying() (*http.Response, error) {
	request, err := http.NewRequest("GET", songUrl, nil)
	if err != nil {
		return nil, err
	}
	return tp.Tunnel.Execute(request)
}

func (tp *Tikplay) Playlist(n uint) (*http.Response, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", queueUrl, n), nil)
	if err != nil {
		return nil, err
	}
	return tp.Tunnel.Execute(request)
}

func (tp *Tikplay) Skip() (*http.Response, error) {
	request, err := http.NewRequest("DELETE", songUrl, nil)
	if err != nil {
		return nil, err
	}
	return tp.Tunnel.Execute(request)
}

func (tp *Tikplay) Clear() (*http.Response, error) {
	request, err := http.NewRequest("DELETE", queueUrl, nil)
	if err != nil {
		return nil, err
	}
	return tp.Tunnel.Execute(request)
}

func (tp *Tikplay) Task(id uint) (*http.Response, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/%d", taskUrl, id), nil)
	if err != nil {
		return nil, err
	}
	return tp.Tunnel.Execute(request)
}

func (tp *Tikplay) Play(url string) (*http.Response, error) {
	jsonbytes, err := json.Marshal(playJSON{tp.Whoami, url})
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", songUrl, bytes.NewReader(jsonbytes))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return tp.Tunnel.Execute(request)
}

func (tp *Tikplay) Close() {
	tp.Tunnel.Close()
}
