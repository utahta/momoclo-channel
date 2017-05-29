package linenotify

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type RequestToken struct {
	Client *http.Client

	Code         string
	RedirectURI  string
	ClientID     string
	ClientSecret string
}

func NewRequestToken(code, redirectURI, clientID, clientSecret string) *RequestToken {
	return &RequestToken{
		Client:       http.DefaultClient,
		Code:         code,
		RedirectURI:  redirectURI,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}
}

func (r *RequestToken) Get() (string, error) {
	v := url.Values{}
	v.Add("grant_type", "authorization_code")
	v.Add("code", r.Code)
	v.Add("redirect_uri", r.RedirectURI)
	v.Add("client_id", r.ClientID)
	v.Add("client_secret", r.ClientSecret)

	resp, err := r.Client.Post(
		"https://notify-bot.line.me/oauth/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(v.Encode()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		var data interface{}
		err = json.NewDecoder(resp.Body).Decode(&data)
		if err != nil {
			return "", err
		}
		root := data.(map[string]interface{})
		return root["access_token"].(string), nil
	}
	return "", errors.New(resp.Status)
}
