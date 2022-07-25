package csrf

import (
	"crypto/sha256"
	"encoding/base64"
	"io"

	"github.com/Edbeer/restapi/pkg/logger"
)

const (
	CSRFHeader = "X-CSRF-Token"
	// 32 bytes
	csrfSalt = "nzB0XYuZJWwm2YfNEz9o1kCQqnP3Bxgl"
)

// Create CSRF token
func MakeToken(sessionId string, logger logger.Logger) string {
	hash := sha256.New()
	_, err := io.WriteString(hash, csrfSalt+sessionId)
	if err != nil {
		logger.Errorf("Make CSRF Token", err)
	}
	token := base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
	return token
}

// Validate CSRF token
func ValidateToken(token string, sid string, logger logger.Logger) bool {
	trueToken := MakeToken(sid, logger)
	return token == trueToken
}
