package key

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

type KeyType int

const (
	RSA KeyType = iota
	ED25519
	ECDSAP384
)

const (
	PRIVATE_BLOCK string = "PRIVATE KEY"
	PUBLIC_BLOCK  string = "PUBLIC KEY"
)

var KeyTypes = map[KeyType][]string{
	RSA:       {"rsa"},
	ED25519:   {"ed25519"},
	ECDSAP384: {"ecdsa-p384"},
}

var KeyTypeDescription = map[KeyType]string{
	RSA:       "implements RSA encryption as specified in PKCS #1 and RFC 8017.",
	ED25519:   "implements the Ed25519 signature algorithm.",
	ECDSAP384: "implements the Elliptic Curve Digital Signature Algorithm, as defined in FIPS 186-4 and SEC 1, Version 2.0.",
}

func generateRSAKey(bitsize uint) (crypto.PrivateKey, crypto.PublicKey, error) {
	if bitsize%256 != 0 {
		return nil, nil, errors.New("bitsize must be a multiple of 256")
	}
	key, err := rsa.GenerateKey(rand.Reader, int(bitsize))
	if err != nil {
		return nil, nil, err
	}
	return key, key.Public(), nil
}

func generateEDKey() (crypto.PrivateKey, crypto.PublicKey, error) {
	return ed25519.GenerateKey(rand.Reader)
}

func generateECDSAKey() (crypto.PrivateKey, crypto.PublicKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return key, key.Public(), nil
}

func MarshalKeyPair(private crypto.PrivateKey, public crypto.PublicKey) (privateBytes, publicBytes []byte, err error) {
	privateBytes, err = x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return nil, nil, err
	}
	publicBytes, err = x509.MarshalPKIXPublicKey(public)
	if err != nil {
		return nil, nil, err
	}

	privateBlock := pem.Block{
		Type:    PRIVATE_BLOCK,
		Bytes:   privateBytes,
		Headers: nil,
	}
	privateBytes = pem.EncodeToMemory(&privateBlock)

	publicBlock := pem.Block{
		Type:    PUBLIC_BLOCK,
		Bytes:   publicBytes,
		Headers: nil,
	}

	privateBytes = append(privateBytes, pem.EncodeToMemory(&publicBlock)...)
	publicBytes = pem.EncodeToMemory(&publicBlock)

	return
}

func findBlock(data []byte, typ string) (*pem.Block, error) {
	var block *pem.Block
	for {
		block, data = pem.Decode(data)
		if block == nil {
			return nil, fmt.Errorf("block '%s' not found", typ)
		}
		if block.Type == typ {
			return block, nil
		}
	}
}

func ParsePrivateKey(keyBytes []byte) (crypto.PrivateKey, error) {
	block, err := findBlock(keyBytes, PRIVATE_BLOCK)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKCS8PrivateKey(block.Bytes)
}

func ParsePublicKey(keyBytes []byte) (crypto.PublicKey, error) {
	block, err := findBlock(keyBytes, PUBLIC_BLOCK)
	if err != nil {
		return nil, err
	}

	return x509.ParsePKIXPublicKey(block.Bytes)
}

func GenerateKeyPair(typ KeyType, bitSize uint) (crypto.PrivateKey, crypto.PublicKey, error) {
	var privateKey crypto.PrivateKey
	var publicKey crypto.PublicKey
	var err error
	switch typ {
	case RSA:
		privateKey, publicKey, err = generateRSAKey(bitSize)
	case ED25519:
		privateKey, publicKey, err = generateEDKey()
	case ECDSAP384:
		privateKey, publicKey, err = generateECDSAKey()
	default:
		err = errors.New("invalid algorithm")
	}
	return privateKey, publicKey, err
}
