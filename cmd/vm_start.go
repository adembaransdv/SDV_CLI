package cmd

import (
    "crypto/tls"
    "fmt"
    "io"
    "net/http"
    vmware "SDVCLI/Auth"
    "github.com/spf13/cobra"
)

var StartCmd = &cobra.Command{
    Use:   "start [vmID]",
    Short: "Démarrer une VM",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        config, err := vmware.CheckConfiguration()
        if err != nil {
            fmt.Println("Aucune configuration trouvée")
            return
        }
        vmID := args[0]
        err = StartVM(config, vmID)
        if err != nil {
            fmt.Println("Impossible de démarrer la VM :", err)
            return
        }
        fmt.Printf("La VM avec l'ID %s a été démarrée avec succès.\n", vmID)
    },
}

func init() {
    ServeurCmd.AddCommand(StartCmd)
}

func StartVM(c *vmware.Client, vmID string) error {
    url := fmt.Sprintf("%s/rest/vcenter/vm/%s/power/start", vmware.Host, vmID)

    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return fmt.Errorf("erreur lors de la création de la requête POST : %v", err)
    }

    req.Header.Add("vmware-api-session-id", c.Token)

    clientHTTP := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    resp, err := clientHTTP.Do(req)
    if err != nil {
        return fmt.Errorf("erreur lors de l'appel HTTP : %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("erreur lors de la lecture de la réponse : %v", err)
    }

    if resp.StatusCode != 200 {
        return fmt.Errorf("erreur lors du démarrage de la VM, code : %d, message : %s", resp.StatusCode, string(body))
    }

    fmt.Println("Réponse de l'API StartVM :", string(body))
    return nil
}
