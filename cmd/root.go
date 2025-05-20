package cmd

import (
    "sdvcli/Auth"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "sdvcli",
    Short: "CLI pour interagir avec vSphere",
}

var VmwareClient *Auth.Client

func Execute() {
    cobra.CheckErr(rootCmd.Execute())
}
