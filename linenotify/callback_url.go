package linenotify

import "github.com/utahta/momoclo-channel/config"

// CallbackURL returns LINE Notify callback URL
func CallbackURL() string {
	return config.C().App.BaseURL + "/line/notify/callback"
}
