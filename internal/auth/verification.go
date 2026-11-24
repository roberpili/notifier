package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/rober0xf/notifier/internal/domain/entity"
)

func GenerateVerificationToken(expirationHours int) (*entity.VerificationToken, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return nil, fmt.Errorf("generate random token: %w", err)
	}

	plainToken := hex.EncodeToString(b)
	hash := sha256.Sum256([]byte(plainToken))
	expiresAt := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	return &entity.VerificationToken{
		Token:     plainToken,
		Hash:      hex.EncodeToString(hash[:]),
		ExpiresAt: expiresAt,
		Timeout:   time.Until(expiresAt),
	}, nil
}
