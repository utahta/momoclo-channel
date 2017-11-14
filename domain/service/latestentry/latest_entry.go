package latestentry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/utahta/momoclo-channel/domain/model"
)

// Parse returns *LatestEntry given url
func Parse(urlStr string) (*model.LatestEntry, error) {
	var code string
	blogCodes := []string{
		model.LatestEntryCodeTamai,
		model.LatestEntryCodeMomota,
		model.LatestEntryCodeAriyasu,
		model.LatestEntryCodeSasaki,
		model.LatestEntryCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(urlStr, fmt.Sprintf("https://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(urlStr, "http://www.tfm.co.jp/clover/") {
		code = model.LatestEntryCodeHappyclo
	} else if strings.HasPrefix(urlStr, "http://www.momoclo.net") {
		code = model.LatestEntryCodeAeNews
	} else if strings.HasPrefix(urlStr, "https://www.youtube.com") {
		code = model.LatestEntryCodeYoutube
	}

	if code == "" {
		return nil, errors.New("code not supported")
	}
	return &model.LatestEntry{ID: code, Code: code, URL: urlStr}, nil
}
