package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	vmware "SDVCLI/Auth"
	database "SDVCLI/Database"

	"github.com/spf13/cobra"
)

var vmName string
var guestOS string

var Createvm = &cobra.Command{
	Use:   "create",
	Short: "To create a new VM",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := vmware.CheckConfiguration()
		if err != nil {
			fmt.Println("Aucune configuration trouvée")
			return
		}
		err = createVM(config, vmware.Host, vmName, guestOS)
		if err != nil {
			log.Fatalf("Erreur lors de la création de la VM : %v", err)
		}
	},
}

func init() {
	Createvm.Flags().StringVarP(&vmName, "name", "n", "", "Nom de la VM")
	Createvm.Flags().StringVarP(&guestOS, "guest-os", "g", "", "Type de guest OS (ex: DEBIAN_10_64)")
	Createvm.MarkFlagRequired("name")
	Createvm.MarkFlagRequired("guest-os")
	ServeurCmd.AddCommand(Createvm)
}

func createVM(client *vmware.Client, host, name, guestOS string) error {

	var result struct {
		Value string `json:"value"`
	}

	payload := map[string]interface{}{
		"spec": map[string]interface{}{
			"name":     name,
			"guest_OS": guestOS,
			"placement": map[string]interface{}{
				"folder":        "group-v22",
				"resource_pool": "resgroup-9",
				"datastore":     "datastore-11",
				"host":          "host-10",
			},
		},
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erreur lors du marshal JSON : %v", err)
	}

	req, err := http.NewRequest("POST", host+"/rest/vcenter/vm", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return errors.New("erreur lors de la création de la requête POST")
	}
	req.Header.Add("vmware-api-session-id", client.Token)
	req.Header.Add("Content-Type", "application/json")
	fmt.Println("Token utilisé pour l'autorisation :", client.Token)

	resp, err := vmware.InsecureHTTPClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de la requête POST :", err)
		return errors.New("Erreur lors de la requete http")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.New("erreur lors de la lecture du corps de la réponse")
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return fmt.Errorf("requête échouée avec le code %d : %s", resp.StatusCode, string(body))
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("erreur lors du parsing JSON : %v", err)
	}
	database.CheckDatabase()
	database.AddKeyToBDD(result.Value, "test")
	fmt.Println("ID de la VM créée :", result.Value)
	return nil
}
