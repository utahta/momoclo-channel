package linenotify

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type RequestNotify struct {
	Client *http.Client
}

var (
	ErrorNotifyInvalidAccessToken = errors.New("Invalid access token.")
)

func NewRequestNotify() *RequestNotify {
	return &RequestNotify{Client: http.DefaultClient}
}

func (r *RequestNotify) Notify(token, message, imageThumbnail, imageFullsize string) error {
	v := url.Values{}
	v.Add("message", message)
	if imageThumbnail != "" {
		v.Add("imageThumbnail", imageThumbnail)
	}
	if imageFullsize != "" {
		v.Add("imageFullsize", imageFullsize)
	}

	req, err := http.NewRequest("POST", "https://notify-api.line.me/api/notify", strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := r.Client.Do(req)
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrorNotifyInvalidAccessToken
	}

	if resp.StatusCode != http.StatusOK {
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		var data interface{}
		err = json.Unmarshal(content, &data)
		if err != nil {
			return err
		}
		root := data.(map[string]interface{})
		return errors.New(root["message"].(string))
	}
	return nil
}
