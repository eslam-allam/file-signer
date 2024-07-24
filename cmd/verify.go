/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"log"

	"github.com/eslam-allam/file-signer/internal/constant"
	"github.com/eslam-allam/file-signer/internal/fs"
	"github.com/eslam-allam/file-signer/internal/key"
	"github.com/eslam-allam/file-signer/internal/licence"
	"github.com/spf13/cobra"
)

var verifyCmdFlags = struct {
	publicKey string
}{}

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify [licence-file]",
	Args:  cobra.ExactArgs(1),
	Short: "Verify a licence file using public key",
	Run: func(cmd *cobra.Command, args []string) {
		signedLicenceBytes, err := fs.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}
		var signedLicence licence.SignedLicence
		err = json.Unmarshal(signedLicenceBytes, &signedLicence)
		if err != nil {
			log.Fatal(err)
		}

		publicBytes, err := fs.ReadFile(verifyCmdFlags.publicKey)
		if err != nil {
			log.Fatal(err)
		}

		publicKey, err := key.ParsePublicKey(publicBytes)
		if err != nil {
			log.Fatal(err)
		}

		err = licence.VerifyLicenceSignature(signedLicence, publicKey)
		if err != nil {
			log.Fatal(err)
		}
		log.Print("Signature valid")
	},
}

func init() {
	licenceCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&verifyCmdFlags.publicKey,
		"public-key", "k", constant.PUBLIC_KEY_FILE_NAME, "Public key used for verifying licence signature")
}
