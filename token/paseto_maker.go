package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

var (
	errInvalidToken = errors.New("invalid token")
)

// PasetoMaker is a PASETO v4-local token maker using a symmetric key.
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker creates a new PasetoMaker with the given symmetric key.
// The key must be at least chacha20poly1305.KeySize (32 bytes) long.
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < chacha20poly1305.KeySize {
		return nil, fmt.Errorf("secret key must be at least %d characters long", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}
	return maker, nil
}

// CreateToken creates a new PASETO token for a specific username and duration.
func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	return paseto.NewV2().Encrypt(maker.symmetricKey, payload, nil)
}

// VerifyToken verifies a PASETO token and returns the payload if it is valid.
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	var payload Payload
	err := maker.paseto.Decrypt(token, maker.symmetricKey, &payload, nil)
	if err != nil {
		return nil, errInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return &payload, nil
}
