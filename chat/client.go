package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:         apiKey,
		Endpoint:       "https://api.apigw.smt.docomo.ne.jp/dialogue/v1",
		userChatStatus: map[string]chatStatus{},
	}
}

type Client struct {
	APIKey         string
	Endpoint       string
	userChatStatus map[string]chatStatus
}

func (c *Client) Talk(userName string, message string) (string, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/dialogue?APIKEY=%s", c.Endpoint, c.APIKey)
	status, _ := c.userChatStatus[userName]
	reqBody := DialogueReq{
		UTT:      message,
		Nickname: userName,
		Context:  status.context,
		Mode:     status.mode,
		T:        "30",
	}
	log.Printf("%#v", reqBody)
	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}
	resp, err := client.Post(uri, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	respBody := &DialogueResp{}
	b = buf.Bytes()
	err = json.Unmarshal(b, respBody)
	if err != nil {
		return "", err
	}

	log.Printf("%#v", string(b))

	status.context = respBody.Context
	status.mode = respBody.Mode

	log.Printf("%#v", status)

	c.userChatStatus[userName] = status
	return respBody.UTT, nil
}

type chatStatus struct {
	context string
	mode    string
}

type DialogueReq struct {
	UTT            string `json:"utt,omitempty"`
	Context        string `json:"context,omitempty"`
	Nickname       string `json:"nickname,omitempty"`
	NicknameY      string `json:"nickname_y,omitempty"`
	Sex            string `json:"sex,omitempty"`
	Bloodtype      string `json:"bloodtype,omitempty"`
	BirthdateY     int    `json:"birthdateY,omitempty"`
	BirthdateM     int    `json:"birthdateM,omitempty"`
	BirthdateD     int    `json:"birthdateD,omitempty"`
	Age            int    `json:"age,omitempty"`
	Constellations string `json:"constellations,omitempty"`
	Place          string `json:"place,omitempty"`
	Mode           string `json:"mode,omitempty"`
	T              string `json:"t,omitempty"` // 20: 関西弁、30赤ちゃんらしい
}

type DialogueResp struct {
	UTT     string `json:"utt,omitempty"`
	Yomi    string `json:"yomi,omitempty"`
	Mode    string `json:"mode,omitempty"`
	DA      string `json:"da,omitempty"`
	Context string `json:"context,omitempty"`
}
