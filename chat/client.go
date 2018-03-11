package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:      apiKey,
		Endpoint:    "https://api.apigw.smt.docomo.ne.jp/dialogue/v1",
		userContext: map[string]string{},
	}
}

type Client struct {
	APIKey      string
	Endpoint    string
	userContext map[string]string
}

func (c *Client) Talk(userName string, message string) (string, error) {
	client := &http.Client{}

	uri := fmt.Sprintf("%s/dialogue?APIKEY=%s", c.Endpoint, c.APIKey)
	context, _ := c.userContext[userName]
	reqBody := DialogueReq{
		UTT:     message,
		Context: context,
	}
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
	err = json.Unmarshal(buf.Bytes(), respBody)
	if err != nil {
		return "", err
	}
	c.userContext[userName] = respBody.Context
	return respBody.UTT, nil
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
}

type DialogueResp struct {
	UTT     string `json:"utt,omitempty"`
	Yomi    string `json:"yomi,omitempty"`
	Mode    string `json:"mode,omitempty"`
	DA      string `json:"da,omitempty"`
	Context string `json:"context,omitempty"`
}
