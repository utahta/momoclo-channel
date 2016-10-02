package model

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/mjibson/goon"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type LineNotification struct {
	Id         string `datastore:"-" goon:"id"`
	TokenCrypt string `datastore:",noindex"`
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

	key := []byte(os.Getenv("LINENOTIFY_TOKEN_KEY"))
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

	key := []byte(os.Getenv("LINENOTIFY_TOKEN_KEY"))
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

//type LineUserQuery struct {
//	context context.Context
//	cursor  datastore.Cursor
//	Limit   int
//}
//
//func NewLineUserQuery(ctx context.Context) *LineUserQuery {
//	return &LineUserQuery{context: ctx, Limit: 100}
//}
//
//func (u *LineUserQuery) GetIds() ([]string, error) {
//	q := datastore.NewQuery("LineUser").Filter("Enabled =", true).KeysOnly().Limit(u.Limit)
//	if u.cursor.String() != "" {
//		q = q.Start(u.cursor)
//	}
//
//	ids := []string{}
//	t := q.Run(u.context)
//	for {
//		k, err := t.Next(nil)
//		if err == datastore.Done {
//			break
//		}
//		if err != nil {
//			return nil, err
//		}
//		ids = append(ids, k.StringID())
//	}
//
//	var err error
//	u.cursor, err = t.Cursor()
//	if err != nil {
//		return nil, err
//	}
//	return ids, nil
//}
