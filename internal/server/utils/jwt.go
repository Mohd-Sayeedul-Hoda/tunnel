package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Mohd-Sayeedul-Hoda/tunnel/internal/server/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenDetails struct {
	Token     string
	TokenUuid string
	UserID    int
	Verified  bool
	ExpiresIn int64
}

var (
	ErrInvalidClaims = errors.New("invalid token claims")
	ErrTokenExpired  = errors.New("token is expired")
)

func CreateToken(user *models.User, ttl time.Duration, privateKey string) (*TokenDetails, error) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: now.Add(ttl).Unix(),
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	td.TokenUuid = uuid.String()
	td.UserID = user.Id
	td.Verified = user.EmailVerified

	cleanPrivateKey := strings.TrimSpace(privateKey)
	decodePrivateKey, err := base64.StdEncoding.DecodeString(cleanPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(decodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = user.Id
	atClaims["token_uuid"] = td.TokenUuid
	atClaims["verified"] = td.Verified
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodEdDSA, atClaims).SignedString(key)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ValidateToken(token string, publicKey string) (*TokenDetails, error) {

	cleanPublicKey := strings.TrimSpace(publicKey)
	decodePublicKey, err := base64.StdEncoding.DecodeString(cleanPublicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decode public key: %w", err)
	}

	key, err := jwt.ParseEdPublicKeyFromPEM(decodePublicKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse token public key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, ErrInvalidClaims
	}

	userIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return &TokenDetails{
		TokenUuid: fmt.Sprint(claims["token_uuid"]),
		Verified:  claims["verified"].(bool),
		UserID:    int(userIDFloat),
	}, nil

}
