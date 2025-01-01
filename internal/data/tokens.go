package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"github.com/kharljhon14/greenlight/internal/validator"
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

func ValidateTokenPlainText(v *validator.Validator, tokenPlaintext string) {
	v.Check(tokenPlaintext != "", "token", "must be provided")
	v.Check(len(tokenPlaintext) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	DB *sql.DB
}

func (m TokenModel) New(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (m TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expirty, scope)
	VALUES ($1, $2, $3)`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, args...)

	return err
}

func (m TokenModel) DeleteAllForUser(scope string, userID int64) error {
	query := `DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, scope, userID)

	return err
}
