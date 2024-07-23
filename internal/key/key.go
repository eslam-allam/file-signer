package key

import (
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
)

type KeyType int

const (
	RSA KeyType = iota
	ED25519
	ECDSA
	ECDH
)

var KeyTypes = map[KeyType][]string{
	RSA:     {"rsa"},
	ED25519: {"ed25519"},
	ECDSA:   {"ecdsa"},
	ECDH:    {"ecdh"},
}

var KeyTypeDescription = map[KeyType]string{
	RSA:     "implements RSA encryption as specified in PKCS #1 and RFC 8017.",
	ED25519: "implements the Ed25519 signature algorithm.",
	ECDSA:   "implements the Elliptic Curve Digital Signature Algorithm, as defined in FIPS 186-4 and SEC 1, Version 2.0.",
	ECDH:    "ecdh implements Elliptic Curve Diffie-Hellman over NIST curves and Curve25519.",
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

func generateECDHKey() (crypto.PrivateKey, crypto.PublicKey, error) {
	key, err := ecdh.P384().GenerateKey(rand.Reader)
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
	return
}

func ParsePrivateKey(keyBytes []byte) (crypto.PrivateKey, error) {
	return x509.ParsePKCS8PrivateKey(keyBytes)
}

func ParsePublicKey(keyBytes []byte) (crypto.PublicKey, error) {
	return x509.ParsePKIXPublicKey(keyBytes)
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
	case ECDSA:
		privateKey, publicKey, err = generateECDSAKey()
	case ECDH:
		privateKey, publicKey, err = generateECDHKey()
	}
	return privateKey, publicKey, err
}
