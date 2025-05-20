package cmd

import (
    "fmt"
    "os"
    "sdvcli/Auth"

    "github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Connexion à vSphere avec session en mémoire",
    Run: func(cmd *cobra.Command, args []string) {
        client, err := Auth.Connection("administrator@vsphere.local", "SDVNantes44!!")
        if err != nil {
            fmt.Println("Login failed:", err)
            os.Exit(1)
        }
        VmwareClient = client
        fmt.Println("Login successful.")
    },
}

func init() {
    rootCmd.AddCommand(loginCmd)
}
