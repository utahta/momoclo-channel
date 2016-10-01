package linenotify

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"
)

const (
	RESPONSE_MODE_FORM_POST = "form_post"
)

type RequestAuthorization struct {
	ClientID     string
	RedirectURI  string
	ResponseMode string
	State        string
}

type ResponseAuthorization struct {
	Code             string
	State            string
	Error            string
	ErrorDescription string
}

func NewRequestAuthorization(clientID, redirectURI string) (*RequestAuthorization, error) {
	state, err := generateState()
	if err != nil {
		return nil, err
	}

	return &RequestAuthorization{
		ClientID:     clientID,
		RedirectURI:  redirectURI,
		ResponseMode: RESPONSE_MODE_FORM_POST,
		State:        state,
	}, nil
}

func (r *RequestAuthorization) AuthorizeURL() (string, error) {
	u, err := url.Parse("https://notify-bot.line.me/oauth/authorize")
	if err != nil {
		return "", err
	}

	v := url.Values{}
	v.Add("response_type", "code")
	v.Add("client_id", r.ClientID)
	v.Add("redirect_uri", r.RedirectURI)
	v.Add("scope", "notify")
	v.Add("state", r.State)
	v.Add("response_mode", r.ResponseMode)
	u.RawQuery = v.Encode()

	return u.String(), nil
}

func (r *RequestAuthorization) Redirect(w http.ResponseWriter, req *http.Request) error {
	url, err := r.AuthorizeURL()
	if err != nil {
		return err
	}
	http.Redirect(w, req, url, http.StatusFound)
	return nil
}

func generateState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func ParseAuthResponse(r *http.Request) *ResponseAuthorization {
	resp := &ResponseAuthorization{}
	resp.Code = r.FormValue("code")
	resp.State = r.FormValue("state")
	resp.Error = r.FormValue("error")
	resp.ErrorDescription = r.FormValue("error_description")
	return resp
}
