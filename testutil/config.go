package testutil

import (
	"io/ioutil"
	"os"

	"github.com/utahta/momoclo-channel/config"
)

// MustConfigLoad loads config file for test
func MustConfigLoad() {
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tmpfile.Name())

	_, err = tmpfile.WriteString(`
[LineNotify]
  TokenKey = "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
`)
	if err != nil {
		panic(err)
	}

	config.MustLoad(tmpfile.Name())
}
