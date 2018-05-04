package config

import (
	"time"

	"github.com/pelletier/go-toml"
	"github.com/utahta/momoclo-channel/timeutil"
)

// Config represents all settings
type Config struct {
	App                App
	Twitter            Twitter
	LineBot            LineBot
	GoogleCustomSearch GoogleCustomSearch
	LineNotify         LineNotify
}

// App represents app entire settings
type App struct {
	BaseURL string
}

// Twitter represents twitter settings
type Twitter struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
	Disabled          bool
}

// LineBot represents LINE Bot settings
type LineBot struct {
	ChannelSecret string
	ChannelToken  string
}

// GoogleCustomSearch represents google custom search api settings
type GoogleCustomSearch struct {
	ApiID  string
	ApiKey string
}

// LineNotify represents LINE Notify settings
type LineNotify struct {
	ClientID     string
	ClientSecret string
	TokenKey     string
	Disabled     bool
}

var (
	c *Config
)

// C returns Config
func C() *Config {
	return c
}

// MustLoad loads config file
// it causes panic if failed
func MustLoad(path string) {
	if err := Load(path); err != nil {
		panic(err)
	}
}

// Load loads config file
func Load(path string) error {
	t, err := toml.LoadFile(path)
	if err != nil {
		return err
	}

	c = &Config{}
	if err := t.Unmarshal(c); err != nil {
		return err
	}

	time.Local = timeutil.JST()
	return nil
}
