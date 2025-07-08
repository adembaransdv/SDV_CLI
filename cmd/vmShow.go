package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	vmware "SDVCLI/Auth"
	database "SDVCLI/Database"

	"github.com/spf13/cobra"
)

var vmShowCmd = &cobra.Command{
	Use:   "show [vm-id]",
	Short: "Affiche les détails d'une VM",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vmID := args[0]
		client, err := vmware.CheckConfiguration()
		if err != nil {
			fmt.Println("Erreur d'authentification :", err)
			return
		}

		url := "https://192.168.1.3/rest/vcenter/vm/" + vmID

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Erreur lors de la création de la requête :", err)
			return
		}

		req.Header.Add("vmware-api-session-id", client.Token)
		resp, err := vmware.InsecureHTTPClient.Do(req)
		if err != nil {
			fmt.Println("Erreur lors de la requête HTTP :", err)
			return
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.StatusCode != 200 {
			if resp.StatusCode == 404 {
				fmt.Println("Le serveur n'existe pas.")
				return
			}
			fmt.Printf("Erreur HTTP %d : %s\n", resp.StatusCode, string(body))
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Erreur lors du parsing JSON :", err)
			return
		}

		database.CheckDatabase()
		valid, err := database.FindInBDD(vmID)
		if valid {

			vm := result["value"].(map[string]interface{})
			fmt.Println("Détails de la VM", vmID)
			fmt.Println("-------------------------")
			fmt.Println("Nom         :", vm["name"])
			fmt.Println("État        :", vm["power_state"])
			fmt.Println("OS invité   :", vm["guest_OS"])
			fmt.Println("CPU         :", vm["cpu_count"], "vCPU")
			fmt.Println("RAM         :", vm["memory_size_MiB"], "MB")
		}
		if !valid {
			fmt.Println("Impossible d'afficher ce que tu souhaites. Bien essayé coco")
		}
	},
}

func init() {
	ServeurCmd.AddCommand(vmShowCmd)
}
