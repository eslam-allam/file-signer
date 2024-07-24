/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/eslam-allam/file-signer/internal/constant"
	"github.com/eslam-allam/file-signer/internal/fs"
	"github.com/eslam-allam/file-signer/internal/key"
	"github.com/eslam-allam/file-signer/internal/licence"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

var (
	keyType         key.KeyType
	bitSize         uint
	targetDirectory string
	overwrite       bool
)

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign [file]",
	Short: "Generate a private/public key pair and sign a file using those keys",
	Args:  cobra.ExactArgs(1),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) >= 1 {
			return []string{}, cobra.ShellCompDirectiveError
		}
		completions, err := fs.ListDirFilter(".", toComplete, []string{".json"})
		if err != nil {
			return []string{}, cobra.ShellCompDirectiveError
		}
		return completions, cobra.ShellCompDirectiveDefault
	},
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		private, public, err := key.GenerateKeyPair(keyType, bitSize)
		if err != nil {
			log.Fatal(err)
		}

		message, err := fs.ReadFile(args[0])
		message = bytes.Trim(message, "\n\t ")
		if err != nil {
			log.Fatal(err)
		}

		var l licence.Licence
		err = json.Unmarshal(message, &l)
		if err != nil {
			log.Fatal(err)
		}

		signed, err := licence.SignLicence(private, l)
		if err != nil {
			log.Fatal(err)
		}

		signedBytes, err := json.MarshalIndent(signed, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		privateBytes, publicBytes, err := key.MarshalKeyPair(private, public)
		if err != nil {
			log.Fatal(err)
		}

		if targetDirectory == "" {
			targetDirectory = filepath.Dir(args[0])
		}

		err = fs.SaveCreateIntermediate(filepath.Join(targetDirectory, constant.SIGNED_LICENCE_FILE_NAME), signedBytes, overwrite)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(filepath.Join(targetDirectory, constant.PRIVATE_KEY_FILE_NAME), privateBytes, overwrite)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(filepath.Join(targetDirectory, constant.PUBLIC_KEY_FILE_NAME), publicBytes, overwrite)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(signCmd)

	te := enumflag.New(
		&keyType,
		"type",
		key.KeyTypes,
		enumflag.EnumCaseInsensitive,
	)
	te.RegisterCompletion(signCmd, "type", key.KeyTypeDescription)

	signCmd.Flags().VarP(te, "type", "t", "Algorithm used for generating private/public key pairs")
	signCmd.Flags().UintVarP(&bitSize, "bit-size", "b", 4096, "Number of bits used if RSA algorithm is used (must be a multiple of 256)")
	signCmd.Flags().StringVarP(&targetDirectory, "target-directory", "d", "", "Directory used to save signed file and key pair. (default $licence_file_directory)")
	signCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing files with generated files")
}
