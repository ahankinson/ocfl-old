package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

var inspectCommand = &cobra.Command{
	Use: "inspect",
	Short: "Inspect an OCFL Object",
	Long: "Inspect an OCFL Object, displaying information about it",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Inspecting an OCFL Object")
	},
}
