package sign

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/asn1"
	"errors"
	"fmt"
	"math/big"
)

func verifyRSASignature(signature, data []byte, key *rsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	err := rsa.VerifyPKCS1v15(key, crypto.SHA256, hashed[:], signature)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}
	return nil
}

func verifyECDSASignature(signature, data []byte, key *ecdsa.PublicKey) error {
	hashed := sha256.Sum256(data)
	var ecdsaSignature struct {
		R, S *big.Int
	}
	_, err := asn1.Unmarshal(signature, &ecdsaSignature)
	if err != nil {
		return fmt.Errorf("error unmarshaling signature: %v", err)
	}
	if !ecdsa.Verify(key, hashed[:], ecdsaSignature.R, ecdsaSignature.S) {
		return fmt.Errorf("signature verification failed")
	}
	return nil
}

func verifyEDSignature(signature, data []byte, key ed25519.PublicKey) error {
	if !ed25519.Verify(key, data, signature) {
		return errors.New("invalid signature")
	}
	return nil
}

func VerifySignature(signature, data []byte, key crypto.PublicKey) error {
	switch key := key.(type) {
	case *rsa.PublicKey:
		return verifyRSASignature(signature, data, key)
	case *ecdsa.PublicKey:
		return verifyECDSASignature(signature, data, key)
	case ed25519.PublicKey:
		return verifyEDSignature(signature, data, key)
	default:
		return errors.New("invalid public key")
	}
}

func SignMessage(key crypto.PrivateKey, data []byte) ([]byte, error) {
	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, errors.New("private key is not a signer")
	}
	hashed := sha256.Sum256(data)
	return signer.Sign(rand.Reader, hashed[:], crypto.SHA256)
}
