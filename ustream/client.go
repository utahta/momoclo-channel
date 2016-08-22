package ustream

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
)

type Client struct {
	HttpClient *http.Client
}

func NewClient() *Client {
	return &Client{HttpClient: http.DefaultClient}
}

func (c *Client) IsLive() (bool, error) {
	resp, err := c.HttpClient.Get("https://api.ustream.tv/channels/4979543.json")
	if err != nil {
		return false, errors.Wrap(err, "Failed to get ustream")
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, errors.Wrap(err, "Failed to read body")
	}

	var data interface{}
	err = json.Unmarshal(content, &data)
	if err != nil {
		return false, errors.Wrap(err, "Failed to unmarshal")
	}
	root := data.(map[string]interface{})
	channel := root["channel"].(map[string]interface{})
	status := channel["status"].(string)

	return status == "live", nil
}
