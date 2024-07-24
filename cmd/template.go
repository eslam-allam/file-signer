/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"path/filepath"

	"github.com/eslam-allam/file-signer/internal/constant"
	"github.com/eslam-allam/file-signer/internal/fs"
	"github.com/eslam-allam/file-signer/internal/licence"
	"github.com/spf13/cobra"
)

var licenceTemplateFlags = struct {
	targetDir string
	overwrite bool
}{}

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Generate a template licence file to be signed",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		licence, schema, err := licence.GetTemplate()
		if err != nil {
			log.Fatal(err)
		}
		err = fs.SaveCreateIntermediate(
			filepath.Join(licenceTemplateFlags.targetDir, constant.LICENCE_FILE_NAME), licence, licenceTemplateFlags.overwrite)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(
			filepath.Join(licenceTemplateFlags.targetDir, constant.SCHEMA_FILE_NAME), schema, licenceTemplateFlags.overwrite)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	licenceCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringVarP(&licenceTemplateFlags.targetDir, "directory", "d", ".", "Target directory where template will be generated")
	templateCmd.Flags().BoolVarP(&licenceTemplateFlags.overwrite, "overwrite", "o", false, "Overwrite template files if they exist")
}
