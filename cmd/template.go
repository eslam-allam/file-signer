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

var (
	templateDir       string
	overwriteTemplate bool
)

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
		err = fs.SaveCreateIntermediate(filepath.Join(templateDir, constant.LICENCE_FILE_NAME), licence, overwriteTemplate)
		if err != nil {
			log.Fatal(err)
		}

		err = fs.SaveCreateIntermediate(filepath.Join(templateDir, constant.SCHEMA_FILE_NAME), schema, overwriteTemplate)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(templateCmd)

	templateCmd.Flags().StringVarP(&templateDir, "directory", "d", ".", "Target directory where template will be generated")
	templateCmd.Flags().BoolVarP(&overwriteTemplate, "overwrite", "o", false, "Overwrite template files if they exist")
}
