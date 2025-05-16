package sfu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type JSEP struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

type ICECandidate struct {
	Candidate     string `json:"candidate"`
	SDPMid        string `json:"sdpMid"`
	SDPMLineIndex uint16 `json:"sdpMLineIndex"`
}

type JanusClient struct {
	BaseURL    string
	HTTPClient *http.Client
	SessionID  int64
	HandleID   int64
}

func NewJanusClient(baseURL string) *JanusClient {
	return &JanusClient{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

func (j *JanusClient) do(body interface{}, url string) ([]byte, error) {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := j.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (j *JanusClient) CreateSession() error {
	url := fmt.Sprintf("%s/janus", j.BaseURL)
	req := map[string]interface{}{
		"janus":       "create",
		"transaction": "txn_create",
	}
	resp, err := j.do(req, url)
	if err != nil {
		return err
	}
	var out struct {
		Data struct{ ID int64 } `json:"data"`
	}
	if err := json.Unmarshal(resp, &out); err != nil {
		return err
	}
	j.SessionID = out.Data.ID
	return nil
}

func (j *JanusClient) Attach(plugin string) error {
	url := fmt.Sprintf("%s/janus/%d", j.BaseURL, j.SessionID)
	req := map[string]interface{}{
		"janus":       "attach",
		"plugin":      plugin,
		"transaction": "txn_attach",
	}
	resp, err := j.do(req, url)
	if err != nil {
		return err
	}
	var out struct {
		Data struct{ ID int64 } `json:"data"`
	}
	if err := json.Unmarshal(resp, &out); err != nil {
		return err
	}
	j.HandleID = out.Data.ID
	return nil
}

func (j *JanusClient) SetupMedia(audio, video bool) (*JSEP, error) {
	url := fmt.Sprintf("%s/janus/%d/%d", j.BaseURL, j.SessionID, j.HandleID)
	req := map[string]interface{}{
		"janus":       "message",
		"body":        map[string]bool{"audio": audio, "video": video},
		"transaction": "txn_setup",
	}
	resp, err := j.do(req, url)
	if err != nil {
		return nil, err
	}
	var out struct {
		Jsep JSEP `json:"jsep"`
	}
	if err := json.Unmarshal(resp, &out); err != nil {
		return nil, err
	}
	return &out.Jsep, nil
}

func (j *JanusClient) SendAnswer(jsep *JSEP) error {
	url := fmt.Sprintf("%s/janus/%d/%d", j.BaseURL, j.SessionID, j.HandleID)
	req := map[string]interface{}{
		"janus":       "message",
		"body":        map[string]string{"request": "start"},
		"transaction": "txn_start",
		"jsep":        jsep,
	}
	_, err := j.do(req, url)
	return err
}

func (j *JanusClient) Tricklet(cand *ICECandidate) error {
	url := fmt.Sprintf("%s/janus/%d/%d", j.BaseURL, j.SessionID, j.HandleID)
	req := map[string]interface{}{
		"janus":       "trickle",
		"candidate":   cand,
		"transaction": "txn_trickle",
	}
	_, err := j.do(req, url)
	return err
}
