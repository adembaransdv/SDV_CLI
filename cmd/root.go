package cmd

import "github.com/spf13/cobra"

var Rootcmd = &cobra.Command{
    Use:   "SDVCLI",
    Short: "CLI",
}

func init() {
    Rootcmd.AddCommand(vmCmd)
}

func Execute() {
    cobra.CheckErr(Rootcmd.Execute())
}
