// root.go — defines the root cobra command and wires all subcommands together.
// The only logic here is registering child commands and calling Execute().
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/mahin/mahin-cli-v1/version"
)

var rootCmd = &cobra.Command{
	Use:   "mahin",
	Short: "Movie CLI — manage your movie & TV show library",
	Long: fmt.Sprintf(`mahin %s — Movie & TV Show Library Manager

A cross-platform CLI tool for managing a personal movie and TV show
library. Scan local folders, clean filenames, fetch metadata from TMDb,
organize files, and track your collection — all from the terminal.

Quick Start:
  mahin movie config set tmdb-key YOUR_API_KEY   Set your TMDb API key
  mahin movie scan ~/Movies                       Scan a folder for videos
  mahin movie ls                                  List your library
  mahin movie search "Inception"                  Search TMDb
  mahin movie info 1                              Show movie details

Management:
  mahin movie move                                Move files interactively
  mahin movie rename                              Batch-rename messy filenames
  mahin movie undo                                Undo last move/rename
  mahin movie play 1                              Play with default player

Discovery:
  mahin movie suggest                             Get recommendations
  mahin movie tag add 1 favorite                  Tag your movies
  mahin movie stats                               Library statistics

System:
  mahin version                                   Show version info
  mahin self-update                               Update via git pull

Documentation: https://github.com/mahin/mahin-cli-v1`, version.Short()),
	Version: version.Short(),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mahin %s\n\n", version.Full())
		cmd.Help()
	},
}

func init() {
	// Keep the CLI surface focused on project commands only.
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetVersionTemplate(fmt.Sprintf("mahin %s\n", version.Full()))
	rootCmd.AddCommand(helloCmd, versionCmd, selfUpdateCmd, movieCmd)
}

// Execute is called by main.go. It is the single public entry point.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
