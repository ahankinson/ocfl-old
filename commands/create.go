package commands

import (
	"fmt"
	"github.com/ocfl/ocfl/libocfl"
	"github.com/spf13/cobra"
)

var createObjectCommand = &cobra.Command{
	Use: "create [flags] [path to object]",
	Short: "Create a new OCFL Object",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		iconf := &libocfl.InventoryConfig{
			DigestAlgorithm: "sha512",
			Id: "",
		}
		vconf := &libocfl.VersionConfig{
			UserName: "Andrew Hankinson",
			UserEmail: "andrew.hankinson@gmail.com",
		}

		ib, _ := libocfl.CreateBlankInventory(iconf, vconf)
		str, _ := ib.ToIndentedJSON()
		s := string(str)

		fmt.Println(s)
	},
}
