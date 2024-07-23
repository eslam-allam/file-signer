package certificate

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
)

type licenceBlockDelimeter string

func (l licenceBlockDelimeter) Line() string {
	var line strings.Builder

	line.WriteString(padding)
	line.WriteString(string(l))
	line.WriteString(padding)
	line.WriteByte('\n')
	return line.String()
}

const (
	padding string = "----"

	sigStart licenceBlockDelimeter = "SIG START"
	sigEnd   licenceBlockDelimeter = "SIG END"

	licenceStart licenceBlockDelimeter = "LICENCE START"
	licenceEnd   licenceBlockDelimeter = "LICENCE END"
)

func signMessage(key crypto.PrivateKey, data []byte) ([]byte, error) {
	signer, ok := key.(crypto.Signer)
	if !ok {
		return nil, errors.New("private key is not a signer")
	}
	hash := sha256.New()
	_, err := hash.Write(data)
	if err != nil {
		return nil, err
	}
	hashed := hash.Sum([]byte{})
	return signer.Sign(rand.Reader, hashed, crypto.SHA256)
}

func SignLicence(key crypto.PrivateKey, data []byte) ([]byte, error) {
	signature, err := signMessage(key, data)
	if err != nil {
		return nil, err
	}
	signatureHex := make([]byte, hex.EncodedLen(len(signature)))
	_ = hex.Encode(signatureHex, signature)
	signatureHex = bytes.ToUpper(signatureHex)

	var result bytes.Buffer

	result.WriteString(licenceStart.Line())
	result.Write(data)
	result.WriteByte('\n')
	result.WriteString(licenceEnd.Line())

	result.WriteByte('\n')

	result.WriteString(sigStart.Line())
	result.Write(signatureHex)
	result.WriteByte('\n')
	result.WriteString(sigEnd.Line())

	return result.Bytes(), nil
}
