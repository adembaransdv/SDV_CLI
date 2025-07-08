package cmd

import (
    "crypto/tls"
    "errors"
    "fmt"
    "io"
    "net/http"
    "strings"

    vmware "SDVCLI/Auth"
    database "SDVCLI/Database"
    "github.com/spf13/cobra"
)

var StopCmd = &cobra.Command{
    Use:   "stop [vmID]",
    Short: "Arrêter une VM",
    Args:  cobra.ExactArgs(1),
    Run: func(cmd *cobra.Command, args []string) {
        config, err := vmware.CheckConfiguration()
        if err != nil {
            fmt.Println("Aucune configuration trouvée")
            return
        }

        vmID := args[0]

        exists, err := database.FindInBDD(vmID)
        if err != nil {
            fmt.Println("Erreur lors de la vérification dans la base :", err)
            return
        }
        if !exists {
            fmt.Printf("La VM avec l'ID %s n'existe pas dans la base.\n", vmID)
            return
        }

        alreadyStopped, err := StopVM(config, vmID)
        if err != nil {
            fmt.Println("Impossible d'arrêter la VM :", err)
            return
        }
        if alreadyStopped {
            fmt.Printf("La VM %s est déjà arrêtée.\n", vmID)
        }
        fmt.Printf("La VM avec l'ID %s a été arrêtée avec succès.\n", vmID)
    },
}

func init() {
    ServeurCmd.AddCommand(StopCmd)
}

func StopVM(c *vmware.Client, vmID string) (bool, error) {
    url := fmt.Sprintf("%s/rest/vcenter/vm/%s/power/stop", vmware.Host, vmID)

    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return false, errors.New("erreur lors de la création de la requête POST")
    }

    req.Header.Add("vmware-api-session-id", c.Token)

    clientHTTP := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }

    resp, err := clientHTTP.Do(req)
    if err != nil {
        return false, fmt.Errorf("erreur lors de l'appel HTTP : %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return false, fmt.Errorf("erreur lors de la lecture de la réponse : %v", err)
    }

    if resp.StatusCode == 200 {
        return false, nil
    }

    if resp.StatusCode == 400 {
        if strings.Contains(string(body), "already_in_desired_state") || strings.Contains(string(body), "already powered off") {
            return true, nil // indique que la VM était déjà arrêtée
        }
    }

    // Toute autre erreur
    return false, fmt.Errorf("erreur lors de l'arrêt de la VM, code : %d, message : %s", resp.StatusCode, string(body))
}
