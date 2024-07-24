package licence

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"path/filepath"
	"time"

	"github.com/eslam-allam/file-signer/internal/constant"
	"github.com/google/uuid"
)

const (
	signatureBlock string = "SIGNATURE"
	licenceBlock   string = "LICENCE"
)

type schemaProperty struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Pattern     string `json:"pattern,omitempty"`
	Format      string `json:"format,omitempty"`
	MinLength   int    `json:"minLength,omitempty"`
}

type licenceSchemaProperties struct {
	LicenceKey schemaProperty `json:"licence_key"`
	Name       schemaProperty `json:"name"`
	Email      schemaProperty `json:"email"`
	Product    schemaProperty `json:"product"`
	Version    schemaProperty `json:"version"`
	Issuer     schemaProperty `json:"issuer"`
	IssueDate  schemaProperty `json:"issue_date"`
	ExpiryDate schemaProperty `json:"expiry_date"`
}

type licenceSchemaDefinition struct {
	Type        string                  `json:"type"`
	Title       string                  `json:"title"`
	Description string                  `json:"description"`
	Properties  licenceSchemaProperties `json:"properties"`
	Required    []string                `json:"required"`
}

var licenceSchema licenceSchemaDefinition = licenceSchemaDefinition{
	Type:        "object",
	Title:       "Licence File",
	Description: "Schema for licence files",
	Properties: licenceSchemaProperties{
		LicenceKey: schemaProperty{
			Type:        "string",
			Description: "Identifier for unique subscription. Will be filled with UUIDV4 if empty.",
			MinLength:   1,
		},
		Name: schemaProperty{
			Type:        "string",
			Description: "Name of purchaser",
			MinLength:   3,
		},
		Email: schemaProperty{
			Type:        "string",
			Description: "Email of purchaser",
			Format:      "email",
		},
		Product: schemaProperty{
			Type:        "string",
			Description: "Name of the applicable product",
			MinLength:   3,
		},
		Version: schemaProperty{
			Type:        "string",
			Description: "Version of licence file. Used by each application differently",
			MinLength:   1,
		},
		Issuer: schemaProperty{
			Type:        "string",
			Description: "Uniquely identifies the issuer of the certificate. Can be an email, Id, etc...",
			MinLength:   3,
		},
		IssueDate: schemaProperty{
			Type:        "string",
			Description: "Date of issuing this licence in this format (yyyy-mm-dd). Will be autofilled if not provided",
			Format:      "date",
		},
		ExpiryDate: schemaProperty{
			Type:        "string",
			Description: "Date of licence expiry in this format (yyyy-mm-dd).",
			Format:      "date",
		},
	},
	Required: []string{"name", "email", "product", "version", "issuer", "expiry_date"},
}

type Licence struct {
	Schema     string `json:"$schema"`
	LicenceKey string `json:"licence_key"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Product    string `json:"product"`
	Version    string `json:"version"`
	Issuer     string `json:"issuer"`
	IssueDate  string `json:"issue_date"`
	ExpiryDate string `json:"expiry_date"`
}

type SignedLicence struct {
	Licence
	Signature string `json:"signature"`
}

func GetTemplate() (licence []byte, schema []byte, err error) {
	licence, err = json.MarshalIndent(Licence{Schema: filepath.Join(".", constant.SCHEMA_FILE_NAME)}, "", "  ")
	if err != nil {
		return nil, nil, err
	}
	schema, err = json.MarshalIndent(licenceSchema, "", "  ")
	return
}

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

func validateLicence(licence Licence) error {
	if licence.Name == "" {
		return errors.New("licence.name cannot be empty")
	}
	if licence.Email == "" {
		return errors.New("licence.email cannot be empty")
	}
	if licence.Issuer == "" {
		return errors.New("licence.issuer cannot be empty")
	}
	if licence.Product == "" {
		return errors.New("licence.product cannot be empty")
	}
	if licence.Version == "" {
		return errors.New("licence.version cannot be empty")
	}
	if licence.ExpiryDate == "" {
		return errors.New("licence.expiry_date cannot be empty")
	}
	return nil
}

func SignLicence(key crypto.PrivateKey, licence Licence) (SignedLicence, error) {
	err := validateLicence(licence)
	if err != nil {
		return SignedLicence{}, err
	}

	if licence.LicenceKey == "" {
		licence.LicenceKey = uuid.New().String()
	}

	if licence.IssueDate == "" {
		licence.IssueDate = time.Now().Format(time.DateOnly)
	}

	licenceData, err := json.Marshal(licence)
	if err != nil {
		return SignedLicence{}, err
	}
	signature, err := signMessage(key, licenceData)
	if err != nil {
		return SignedLicence{}, err
	}
	encodedSignature := base64.StdEncoding.EncodeToString(signature)
	return SignedLicence{Licence: licence, Signature: encodedSignature}, nil
}
