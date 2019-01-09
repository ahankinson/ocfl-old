package commands

import (
	"github.com/spf13/cobra"
)

var versionCommand = &cobra.Command{
	Use: "version",
	Short: "Work with OCFL Object Versions",
	Long: "A longer description of what this does",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var stageVersionCommand = &cobra.Command{
	Use: "stage",
	Short: "Stage files to be added to OCFL Object",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var commitVersionCommand = &cobra.Command{
	Use: "commit",
	Short: "Commit staged files to the OCFL Object",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

var statusVersionCommand = &cobra.Command{
	Use: "status",
	Short: "View status of the OCFL Object (staged/committed files)",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
