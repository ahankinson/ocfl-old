package commands

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var rootCmd = &cobra.Command{
	Use:              "ocfl",
	Short:            "An OCFL Client",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logFmt, _ := cmd.Flags().GetString("logFormat")

		if logFmt == "console" {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
		}

		logLevel, _ := cmd.Flags().GetString("logLevel")

		// Set the loglevel to ERROR by default
		if logLevel == "INFO" {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		} else if logLevel == "DEBUG" {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else if logLevel == "WARN" {
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		}
	},
}

func Execute(args []string) {
	rootCmd.AddCommand(validateCommand, createObjectCommand, inspectCommand, versionCommand)
	rootCmd.PersistentFlags().StringP("logFormat", "f", "console", "Output format (console, json)")
	rootCmd.PersistentFlags().StringP("logLevel", "g", "ERROR", "Log Level (DEBUG, INFO, WARN, ERROR)")

	createObjectCommand.Flags().StringP("algorithm", "a", "sha512", "Algorithm choice (md5, sha1, sha256, sha512, blake2b)")
	createObjectCommand.Flags().StringP("creatorEmail", "m", "nobody@example.com", "E-mail address of creator")
	createObjectCommand.Flags().StringP("creatorName", "n", "Jill Blogs", "Name of creator")
	createObjectCommand.Flags().StringP("id", "i", "", "ID of object (default is random UUID)")

	versionCommand.AddCommand(stageVersionCommand, commitVersionCommand, statusVersionCommand)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
