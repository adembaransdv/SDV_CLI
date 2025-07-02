package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	vmware "SDVCLI/Auth"
	database "SDVCLI/Database"

	"github.com/spf13/cobra"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete [vmID]",
	Short: "Supprimer une VM",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		vmID := args[0]

		config, err := vmware.CheckConfiguration()
		if err != nil {
			fmt.Println("Aucune configuration trouvée :", err)
			os.Exit(1)
		}

		exists, err := database.FindInBDD(vmID)
		if err != nil {
			fmt.Println("Erreur lors de la vérification dans la base :", err)
			os.Exit(1)
		}
		if !exists {
			fmt.Printf("La VM avec l'ID %s n'existe pas dans la base.\n", vmID)
			os.Exit(1)
		}

		err = DeleteVM(config, vmID)
		if err != nil {
			fmt.Printf("Erreur lors de la suppression de la VM : %v\n", err)
			os.Exit(1)
		}

		err = database.DeleteFromBDD(vmID)
		if err != nil {
			fmt.Printf("Erreur lors de la suppression dans la base : %v\n", err)
			os.Exit(1)
		}

		fmt.Println("VM supprimée avec succès, entrée supprimée dans la base.")
	},
}

func init() {
	ServeurCmd.AddCommand(DeleteCmd)
}

func DeleteVM(c *vmware.Client, vmID string) error {
	url := vmware.Host + "/rest/vcenter/vm/" + vmID

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.New("erreur lors de la création de la requête DELETE")
	}

	req.Header.Add("vmware-api-session-id", c.Token)

	resp, err := vmware.InsecureHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("Erreur lors de la requête DELETE : %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 && resp.StatusCode != 204 {
		return fmt.Errorf("échec suppression VM, status %d : %s", resp.StatusCode, string(body))
	}

	return nil
}
