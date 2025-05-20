
package cmd

import (
    "fmt"
    "io"
    "net/http"

    "github.com/spf13/cobra"
)

var vmID string

var vmShowCmd = &cobra.Command{
    Use:   "vm show",
    Short: "Afficher les infos d'une VM",
    Run: func(cmd *cobra.Command, args []string) {
        if VmwareClient == nil {
            fmt.Println("Erreur : vous devez d'abord vous connecter avec 'login'")
            return
        }

        url := "https://192.168.1.3/rest/vcenter/vm/" + vmID
        req, _ := http.NewRequest("GET", url, nil)
        req.Header.Set("vmware-api-session-id", VmwareClient.Token)

        client := &http.Client{}
        resp, _ := client.Do(req)
        defer resp.Body.Close()

        body, _ := io.ReadAll(resp.Body)
        fmt.Println(string(body))
    },
}

func init() {
    vmShowCmd.Flags().StringVar(&vmID, "id", "", "ID de la VM")
    rootCmd.AddCommand(vmShowCmd)
}
