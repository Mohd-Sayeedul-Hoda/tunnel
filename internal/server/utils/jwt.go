package utils

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenDetails struct {
	Token     string
	TokenUuid string
	UserID    int
	ExpiresIn int64
}

func CreateToken(userId int, ttl time.Duration, privateKey string) (*TokenDetails, error) {
	now := time.Now().UTC()
	td := &TokenDetails{
		ExpiresIn: now.Add(ttl).Unix(),
	}

	uuid, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}

	td.TokenUuid = uuid.String()
	td.UserID = userId

	decodePrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to decode private key: %w", err)
	}
	key, err := jwt.ParseEdPrivateKeyFromPEM(decodePrivateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse token private key: %w", err)
	}

	atClaims := make(jwt.MapClaims)
	atClaims["sub"] = userId
	atClaims["token_uuid"] = td.TokenUuid
	atClaims["exp"] = td.ExpiresIn
	atClaims["iat"] = now.Unix()
	atClaims["nbf"] = now.Unix()

	td.Token, err = jwt.NewWithClaims(jwt.SigningMethodEdDSA, atClaims).SignedString(key)
	if err != nil {
		return nil, err
	}

	return td, nil
}

func ValidetToken(token string, publicKey string) (*TokenDetails, error) {
	decodePublicKey, err := base64.StdEncoding.DecodeString(publicKey)
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
		return nil, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return &TokenDetails{
		TokenUuid: fmt.Sprint(claims["token_uuid"]),
		UserID:    claims["sub"].(int),
	}, nil

}
