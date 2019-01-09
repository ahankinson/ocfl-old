package commands

import (
	"github.com/ocfl/ocfl/libocfl"
	"github.com/spf13/cobra"
)

var validateCommand = &cobra.Command{
	Use: "validate [flags] [path to validate]",
	Short: "Validate an OCFL Object",
	Long: "A longer description of what this does",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		obj := libocfl.Open(args[0])
		libocfl.Validate(obj)
	},
}
