package main

import (
	"log"
	"os"
	"encoding/json"
	"io/ioutil"

	"golang.org/x/oauth2"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"google.golang.org/api/taskqueue/v1beta2"
)

func main() {
	err := doMain()
	if err != nil {
		log.Fatal(err)
	}
}

func doMain() error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	s, err := newService()
	if err != nil {
		return err
	}
	q, err := s.Taskqueues.Get("momoclo-channel", "queue-line").Do()
	if err != nil {
		return err
	}
	log.Printf("queue:%#v", q)

	return nil
}

func newService() (*taskqueue.Service, error) {
	ts, err := retrieveTokenSource()
	if err != nil {
		return nil, err
	}

	cli := oauth2.NewClient(oauth2.NoContext, ts)
	s, err := taskqueue.New(cli)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func retrieveTokenSource() (*oauth2.TokenSource, error) {
	tf := tokenFile{ FileName: ".oauth2_token" }
	token, err := tf.Load()
	if err != nil {
		return nil, err
	}
	if token == nil {
		token = &oauth2.Token{ RefreshToken: os.Getenv("OAUTH2_REFRESH_TOKEN") }
	}

	conf := &oauth2.Config{
		ClientID:     os.Getenv("OAUTH2_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH2_CLIENT_SECRET"),
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	ts := conf.TokenSource(oauth2.NoContext, token)
	token, err = ts.Token()
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get token.")
	}
	if err := tf.Save(token); err != nil {
		return nil, err
	}
	return ts, nil
}

type tokenFile struct {
	FileName string
}

func (t *tokenFile) Load() (*oauth2.Token, error) {
	if _, err := os.Stat(t.FileName); err != nil {
		// probably file not found.
		return nil, nil
	}

	b, err := ioutil.ReadFile(t.FileName)
	if err != nil {
		return nil, err
	}

	token := &oauth2.Token{}
	if err := json.Unmarshal(b, token); err != nil {
		return nil, errors.Wrap(err, "Failed to unmarshal")
	}
	return token, nil
}

func (t *tokenFile) Save(token *oauth2.Token) error {
	b, err := json.Marshal(token)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(t.FileName, b, 0600); err != nil {
		return err
	}
	return nil
}
