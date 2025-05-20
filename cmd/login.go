package cmd

import (
    "fmt"
    "os"
    "sdvcli/vmware"

    "github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Login to vSphere and store session in memory",
    Run: func(cmd *cobra.Command, args []string) {
        client, err := vmware.CheckConfiguration()
        if err != nil {
            fmt.Println("Login failed:", err)
            os.Exit(1)
        }
        fmt.Printf("Login successful. Token: %s\n", client.Token)
    },
}
