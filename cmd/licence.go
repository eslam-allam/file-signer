/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// licenceCmd represents the licence command
var licenceCmd = &cobra.Command{
	Use:   "licence",
	Short: "Do various licence file operations",
}

func init() {
	rootCmd.AddCommand(licenceCmd)
}
