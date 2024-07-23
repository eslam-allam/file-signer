/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"log"
	"path/filepath"

	"github.com/eslam-allam/file-signer/internal/certificate"
	"github.com/eslam-allam/file-signer/internal/fs"
	"github.com/eslam-allam/file-signer/internal/key"
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
		completions, err := fs.ListDirStartsWithNotEndsWith(".", toComplete, []string{".pem"})
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

		signed, err := certificate.SignLicence(private, message)
		if err != nil {
			log.Fatal(err)
		}

		privateBytes, publicBytes, err := key.MarshalKeyPair(private, public)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(filepath.Join(targetDirectory, "signed.txt"), signed, overwrite)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(filepath.Join(targetDirectory, "private.pem"), privateBytes, overwrite)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(filepath.Join(targetDirectory, "public.pem"), publicBytes, overwrite)
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
	signCmd.Flags().StringVarP(&targetDirectory, "target-directory", "d", ".", "Directory used to save signed file and key pair")
	signCmd.Flags().BoolVarP(&overwrite, "overwrite", "o", false, "Overwrite existing files with generated files")
}
