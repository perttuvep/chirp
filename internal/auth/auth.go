package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPass(pw string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPwHash(pw, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw))
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	issuedat := jwt.NewNumericDate(time.Now().UTC())
	expiresat := jwt.NewNumericDate(time.Now().UTC().Add(expiresIn))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{Issuer: "chirpy", IssuedAt: issuedat, ExpiresAt: expiresat, Subject: userID.String()})
	ss, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return ss, nil
}

// ValidateJWT validates a JWT token and returns the associated user ID.
// It returns the user ID as an uuid.UUID and an error if validation fails.
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	stringId, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(stringId)

	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func GetBearerToken(header http.Header) (string, error) {
	// Get the Authorization header
	auth := header.Get("Authorization")
	if auth == "" {
		return "", errors.New("Authorization header not found")
	}

	// Split by space
	parts := strings.Fields(auth)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("Invalid Authorization header format")
	}

	// The token is the second part
	return parts[1], nil

}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32) // 32 bytes * 8 bits/byte = 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func GetAPIKey(headers http.Header) (string, error) {
	// Get the Authorization header
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", errors.New("ApiKey header not found")
	}

	// Split by space
	parts := strings.Fields(auth)
	if len(parts) != 2 || parts[0] != "ApiKey" {
		return "", errors.New("Invalid ApiKey header format")
	}

	// The token is the second part
	return parts[1], nil

}
