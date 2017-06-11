package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"github.com/utahta/momoclo-channel/lib/config"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type LineNotification struct {
	Id         string `datastore:"-" goon:"id"`
	TokenCrypt string `datastore:",noindex"`
	Admin      bool
	CreatedAt  time.Time
}

func NewLineNotification(token string) (*LineNotification, error) {
	tokenBytes := []byte(token)
	tokenSum := sha256.Sum256(tokenBytes)
	tokenHash := hex.EncodeToString(tokenSum[:])

	// Encryption
	cipherText := make([]byte, aes.BlockSize+len(tokenBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	key := []byte(config.C.Linenotify.TokenKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCTR(block, iv)
	mode.XORKeyStream(cipherText[aes.BlockSize:], tokenBytes)
	tokenCrypt := hex.EncodeToString(cipherText)

	return &LineNotification{Id: tokenHash, TokenCrypt: tokenCrypt}, nil
}

func (l *LineNotification) Token() (string, error) {
	// Decryption
	cipherText, err := hex.DecodeString(l.TokenCrypt)
	if err != nil {
		return "", err
	}

	key := []byte(config.C.Linenotify.TokenKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	token := make([]byte, len(cipherText[aes.BlockSize:]))
	mode := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	mode.XORKeyStream(token, cipherText[aes.BlockSize:])
	return string(token), nil
}

func (l *LineNotification) Put(ctx context.Context) error {
	g := goon.FromContext(ctx)

	// check for cached item
	err := g.Get(l)
	if err == nil {
		return errors.Errorf("LineNotification already exists.")
	} else if err != datastore.ErrNoSuchEntity {
		return err
	}

	return g.RunInTransaction(func(g *goon.Goon) error {
		err := g.Get(l)
		if err != datastore.ErrNoSuchEntity {
			return err
		}

		jst, err := time.LoadLocation("Asia/Tokyo")
		if err != nil {
			return err
		}
		l.CreatedAt = time.Now().In(jst)

		_, err = g.Put(l)
		return err
	}, nil)
}

func (l *LineNotification) Delete(ctx context.Context) error {
	g := goon.FromContext(ctx)
	return g.Delete(g.Key(l))
}

type LineNotificationQuery struct {
	context context.Context
}

func NewLineNotificationQuery(ctx context.Context) *LineNotificationQuery {
	return &LineNotificationQuery{context: ctx}
}

func (l *LineNotificationQuery) GetAll() ([]*LineNotification, error) {
	g := goon.FromContext(l.context)

	items := []*LineNotification{}
	q := datastore.NewQuery("LineNotification")

	_, err := g.GetAll(q, &items)
	if err != nil {
		return nil, err
	}
	return items, nil
}
