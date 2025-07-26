package utils

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"testing"
	"time"
)

func generateTestKeys(t *testing.T) (privateKeyB64 string, publicKeyB64 string) {
	// Generate a new ed25519 key pair.
	pubKey, privKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("Failed to generate ed25519 key pair: %v", err)
	}

	// Convert the private key to the PKCS8 format.
	pkcs8PrivateKey, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		t.Fatalf("Failed to marshal private key to PKCS8: %v", err)
	}

	// Create a PEM block for the private key.
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8PrivateKey,
	})

	// Convert the public key to the PKIX format.
	pkixPublicKey, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		t.Fatalf("Failed to marshal public key to PKIX: %v", err)
	}

	// Create a PEM block for the public key.
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pkixPublicKey,
	})

	// Return the base64 encoded PEM strings.
	return base64.StdEncoding.EncodeToString(privateKeyPEM), base64.StdEncoding.EncodeToString(publicKeyPEM)
}

func TestCreateAndValidateToken(t *testing.T) {
	// Generate keys for the test.
	privateKey, publicKey := generateTestKeys(t)

	userID := "test-user-123"
	ttl := 5 * time.Minute

	createdTokenDetails, err := CreateToken(userID, ttl, privateKey)
	if err != nil {
		t.Fatalf("CreateToken() returned an unexpected error: %v", err)
	}
	if createdTokenDetails == nil {
		t.Fatalf("CreateToken() returned nil token details")
	}
	if createdTokenDetails.Token == "" {
		t.Fatalf("CreateToken() returned an empty token string")
	}

	validatedTokenDetails, err := ValidetToken(createdTokenDetails.Token, publicKey)
	if err != nil {
		t.Fatalf("ValidetToken() returned an unexpected error: %v", err)
	}
	if validatedTokenDetails == nil {
		t.Fatalf("ValidetToken() returned nil token details")
	}

	if validatedTokenDetails.UserID != userID {
		t.Errorf("Expected UserID to be %q, but got %q", userID, validatedTokenDetails.UserID)
	}

	if validatedTokenDetails.TokenUuid != createdTokenDetails.TokenUuid {
		t.Errorf("Expected TokenUuid to be %q, but got %q", createdTokenDetails.TokenUuid, validatedTokenDetails.TokenUuid)
	}
}
