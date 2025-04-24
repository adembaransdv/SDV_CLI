package cmd

import (
    "github.com/spf13/cobra"
)

var vmCmd = &cobra.Command{
    Use:   "vm",
    Short: "Commandes liées aux VMs",
}

func init() {
    vmCmd.AddCommand(StopCmd) 
}
