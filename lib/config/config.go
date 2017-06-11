package config

import (
	"github.com/pelletier/go-toml"
)

type Config struct {
	Twitter            Twitter
	Linebot            Linebot
	GoogleCustomSearch GoogleCustomSearch
	Linenotify         Linenotify
}

type Twitter struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	Disabled          bool
}

type Linebot struct {
	ChannelSecret string
	ChannelToken  string
}

type GoogleCustomSearch struct {
	ApiID  string
	ApiKey string
}

type Linenotify struct {
	ClientID     string
	ClientSecret string
	TokenKey     string
	Disabled     bool
}

var C Config

func MustLoad(path string) {
	if err := Load(path); err != nil {
		panic(err)
	}
}

func Load(path string) error {
	t, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

	if err := t.Unmarshal(&C); err != nil {
		return err
	}
	return nil
}
