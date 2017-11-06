package entity

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/appengine/datastore"
)

const (
	LatestEntryCodeTamai    = "tamai-sd"
	LatestEntryCodeMomota   = "momota-sd"
	LatestEntryCodeAriyasu  = "ariyasu-sd"
	LatestEntryCodeSasaki   = "sasaki-sd"
	LatestEntryCodeTakagi   = "takagi-sd"
	LatestEntryCodeHappyclo = "happyclo"
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

func ParseLatestEntry(urlStr string) (*LatestEntry, error) {
	code := ParseLatestEntryCode(urlStr)
	if code == "" {
		// not supported
		return nil, errors.New("code not supported")
	}
	return &LatestEntry{ID: code, Code: code, URL: urlStr}, nil
}

func ParseLatestEntryCode(urlStr string) string {
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
	}
	return code
}
