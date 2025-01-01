package data

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	// Token instance
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	// Init a zero valued byte slice with a length of 16 bytes
	randomBytes := make([]byte, 16)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}

	// Base-32 strings may be padded at the end with the =
	// character. We don't need this padding character for the purpose of our tokens.
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Generate SHA-256
	hash := sha256.Sum256([]byte(token.Plaintext))
	// Convert it to a slice using the [:] operator before storing it.
	token.Hash = hash[:]

	return token, nil
}
