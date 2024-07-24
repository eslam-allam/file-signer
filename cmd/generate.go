/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"path/filepath"

	"github.com/eslam-allam/file-signer/internal/constant"
	"github.com/eslam-allam/file-signer/internal/fs"
	"github.com/eslam-allam/file-signer/internal/key"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

var generateKeyFlags = struct {
	targetDirectory string
	keyType         key.KeyType
	bitsize         uint
	overwrite       bool
}{}

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate [directory]",
	Short: "Generate a private/public keypair and save them in [directory]",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		privateKey, publicKey, err := key.GenerateKeyPair(generateKeyFlags.keyType, generateKeyFlags.bitsize)
		if err != nil {
			log.Fatal(err)
		}
		privateBytes, publicBytes, err := key.MarshalKeyPair(privateKey, publicKey)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(
			filepath.Join(generateKeyFlags.targetDirectory, constant.PRIVATE_KEY_FILE_NAME), privateBytes, generateKeyFlags.overwrite)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(
			filepath.Join(generateKeyFlags.targetDirectory, constant.PUBLIC_KEY_FILE_NAME), publicBytes, generateKeyFlags.overwrite)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	keyCmd.AddCommand(generateCmd)

	te := enumflag.New(
		&generateKeyFlags.keyType,
		"type",
		key.KeyTypes,
		enumflag.EnumCaseInsensitive,
	)
	te.RegisterCompletion(generateCmd, "type", key.KeyTypeDescription)

	generateCmd.Flags().StringVarP(&generateKeyFlags.targetDirectory, "target-directory", "d", ".", "Directory used to save generated key pair")
	generateCmd.Flags().BoolVarP(&generateKeyFlags.overwrite, "overwrite", "o", false, "Overwrite existing files with generated files")
	generateCmd.Flags().VarP(te, "type", "t", "Algorithm used for generating private/public key pairs")
	generateCmd.Flags().UintVarP(&generateKeyFlags.bitsize, "bit-size", "b", 4096, "Number of bits used if RSA algorithm is used (must be a multiple of 256)")
}
