package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

const (
	// LatestEntryCode* defines identify code
	LatestEntryCodeTamai    = "tamai-sd"
	LatestEntryCodeMomota   = "momota-sd"
	LatestEntryCodeAriyasu  = "ariyasu-sd"
	LatestEntryCodeSasaki   = "sasaki-sd"
	LatestEntryCodeTakagi   = "takagi-sd"
	LatestEntryCodeHappyclo = "happyclo"
	LatestEntryCodeAeNews   = "aenews"
	LatestEntryCodeYoutube  = "youtube"
)

type (
	// LatestEntry for confirm last updated entry url
	LatestEntry struct {
		ID        string `datastore:"-" goon:"id"`
		Code      string
		URL       string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	// LatestEntryRepository interface
	LatestEntryRepository interface {
		Save(*LatestEntry) error
		GetTamaiURL() string
		GetMomotaURL() string
		GetAriyasuURL() string
		GetSasakiURL() string
		GetTakagiURL() string
		GetHappycloURL() string
	}
)

func (l *LatestEntry) SetCreatedAt(t time.Time) {
	l.CreatedAt = t
}

func (l *LatestEntry) GetCreatedAt() time.Time {
	return l.CreatedAt
}

func (l *LatestEntry) SetUpdatedAt(t time.Time) {
	l.UpdatedAt = t
}

func (l *LatestEntry) Load(p []datastore.Property) error {
	return load(l, p)
}

func (l *LatestEntry) Save() ([]datastore.Property, error) {
	return save(l)
}

// ParseLatestEntry creates LatestEntry given url
func ParseLatestEntry(urlStr string) (*LatestEntry, error) {
	code, err := ParseLatestEntryCode(urlStr)
	if err != nil {
		return nil, err
	}
	return &LatestEntry{ID: code, Code: code, URL: urlStr}, nil
}

// ParseLatestEntryCode gets identify code given url
func ParseLatestEntryCode(urlStr string) (string, error) {
	var code string
	blogCodes := []string{
		LatestEntryCodeTamai,
		LatestEntryCodeMomota,
		LatestEntryCodeAriyasu,
		LatestEntryCodeSasaki,
		LatestEntryCodeTakagi,
	}
	for _, c := range blogCodes {
		if strings.HasPrefix(urlStr, fmt.Sprintf("https://ameblo.jp/%s", c)) {
			code = c
			break
		}
	}
	if strings.HasPrefix(urlStr, "http://www.tfm.co.jp/clover/") {
		code = LatestEntryCodeHappyclo
	} else if strings.HasPrefix(urlStr, "http://www.momoclo.net") {
		code = LatestEntryCodeAeNews
	} else if strings.HasPrefix(urlStr, "https://www.youtube.com") {
		code = LatestEntryCodeYoutube
	}

	if code == "" {
		return "", errors.New("code not supported")
	}
	return code, nil
}
