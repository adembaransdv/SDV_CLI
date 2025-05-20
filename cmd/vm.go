package cmd

import (
	"github.com/spf13/cobra"
)

var ServeurCmd = &cobra.Command{
	Use:   "vm",
	Short: "Commandes liées au vm",
}

func init() {
	Rootcmd.AddCommand(ServeurCmd)
}
