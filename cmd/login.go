
package cmd

import (
    "fmt"
    "sdvcli/vmware"

    "github.com/spf13/cobra"
)

var VmwareClient *vmware.Client

var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Connexion à vSphere",
    Run: func(cmd *cobra.Command, args []string) {
        c, err := vmware.Connection("administrator@vsphere.local", "SDVNantes44!!")
        if err != nil {
            fmt.Println("Erreur de connexion:", err)
            return
        }
        VmwareClient = c
        fmt.Println("Connexion réussie.")
    },
}

func init() {
    rootCmd.AddCommand(loginCmd)
}
