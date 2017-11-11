package latest_entry

import (
	"errors"
	"fmt"
	"strings"

	"github.com/utahta/momoclo-channel/domain/entity"
)

// Parse returns *LatestEntry given url
func Parse(urlStr string) (*entity.LatestEntry, error) {
	var code string
	blogCodes := []string{
		entity.LatestEntryCodeTamai,
		entity.LatestEntryCodeMomota,
		entity.LatestEntryCodeAriyasu,
		entity.LatestEntryCodeSasaki,
		entity.LatestEntryCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(urlStr, fmt.Sprintf("https://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(urlStr, "http://www.tfm.co.jp/clover/") {
		code = entity.LatestEntryCodeHappyclo
	} else if strings.HasPrefix(urlStr, "http://www.momoclo.net") {
		code = entity.LatestEntryCodeAeNews
	} else if strings.HasPrefix(urlStr, "https://www.youtube.com") {
		code = entity.LatestEntryCodeYoutube
	}

	if code == "" {
		return nil, errors.New("code not supported")
	}
	return &entity.LatestEntry{ID: code, Code: code, URL: urlStr}, nil
}
