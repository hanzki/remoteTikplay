package tikplay

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
)

type PlayJSON struct {
	User string `json:"user"`
	Url  string `json:"url"`
}

func NowPlaying() (*http.Request, error) {
	return http.NewRequest("GET", songUrl, nil)
}

func Playlist(n uint) (*http.Request, error) {
	return http.NewRequest("GET", fmt.Sprintf("%s/%d", queueUrl, n), nil)
}

func Skip() (*http.Request, error) {
	return http.NewRequest("DELETE", songUrl, nil)
}

func Clear() (*http.Request, error) {
	return http.NewRequest("DELETE", queueUrl, nil)
}

func Play(pjson *PlayJSON) (*http.Request, error) {
	jsonbytes, err := json.Marshal(pjson)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", songUrl, bytes.NewReader(jsonbytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
