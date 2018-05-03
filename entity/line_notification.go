package entity

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"
)

type (
	// LineNotification represents user tokens that published by LINE Notify
	// the token is encrypted
	LineNotification struct {
		ID         string `datastore:"-" goon:"id" validate:"required"`
		TokenCrypt string `datastore:",noindex" validate:"required"`
		Admin      bool
		CreatedAt  time.Time `validate:"required"`
	}
)

// NewLineNotification returns LineNotification given key and token
func NewLineNotification(tokenKey, token string) (*LineNotification, error) {
	tokenBytes := []byte(token)
	tokenSum := sha256.Sum256(tokenBytes)
	tokenHash := hex.EncodeToString(tokenSum[:])

	// Encrypt
	cipherText := make([]byte, aes.BlockSize+len(tokenBytes))
	iv := cipherText[:aes.BlockSize]
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	key := []byte(tokenKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCTR(block, iv)
	mode.XORKeyStream(cipherText[aes.BlockSize:], tokenBytes)
	tokenCrypt := hex.EncodeToString(cipherText)

	return &LineNotification{ID: tokenHash, TokenCrypt: tokenCrypt}, nil
}

// Token returns decrypted token
func (l *LineNotification) Token(tokenKey string) (string, error) {
	// Decrypt
	cipherText, err := hex.DecodeString(l.TokenCrypt)
	if err != nil {
		return "", err
	}

	key := []byte(tokenKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	token := make([]byte, len(cipherText[aes.BlockSize:]))
	mode := cipher.NewCTR(block, cipherText[:aes.BlockSize])
	mode.XORKeyStream(token, cipherText[aes.BlockSize:])
	return string(token), nil
}

// SetCreatedAt sets given time to CreatedAt
func (l *LineNotification) SetCreatedAt(t time.Time) {
	l.CreatedAt = t
}

// GetCreatedAt gets CreatedAt
func (l *LineNotification) GetCreatedAt() time.Time {
	return l.CreatedAt
}

// BeforeSave hook
func (l *LineNotification) BeforeSave() {
	beforeSave(l)
}
