/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/eslam-allam/file-signer/internal/constant"
	"github.com/eslam-allam/file-signer/internal/fs"
	"github.com/eslam-allam/file-signer/internal/key"
	"github.com/eslam-allam/file-signer/internal/licence"
	"github.com/spf13/cobra"
)

var signCmdFlags = struct {
	privateKey      string
	targetDirectory string
	overwrite       bool
}{}

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign [file]",
	Short: "Sign a licence file using the specified private key",
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
	Run: func(cmd *cobra.Command, args []string) {
		privateKey, err := fs.ReadFile(signCmdFlags.privateKey)
		if err != nil {
			log.Fatal(err)
		}

		private, err := key.ParsePrivateKey(privateKey)
		if err != nil {
			log.Fatal(err)
		}

		message, err := fs.ReadFile(args[0])
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

		if signCmdFlags.targetDirectory == "" {
			signCmdFlags.targetDirectory = filepath.Dir(args[0])
		}

		err = fs.SaveCreateIntermediate(filepath.Join(
			signCmdFlags.targetDirectory, constant.SIGNED_LICENCE_FILE_NAME),
			signedBytes, signCmdFlags.overwrite)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	licenceCmd.AddCommand(signCmd)

	signCmd.Flags().StringVarP(&signCmdFlags.privateKey, "private-key", "k", constant.PRIVATE_KEY_FILE_NAME, "Private key used to sign the licence")
	signCmd.Flags().StringVarP(&signCmdFlags.targetDirectory, "target-directory", "d", "", "Directory used to save signed file. (default $licence_file_directory)")
	signCmd.Flags().BoolVarP(&signCmdFlags.overwrite, "overwrite", "o", false, "Overwrite existing files with generated files")
}
