package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	vmware "SDVCLI/Auth"

	"github.com/spf13/cobra"
)

type VM struct {
	VM         string `json:"vm"`
	Name       string `json:"name"`
	PowerState string `json:"power_state"`
}

type VMResponse struct {
	Value []VM `json:"value"`
}

var StartCmd = &cobra.Command{
	Use:   "list",
	Short: "Liste les vm",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := vmware.CheckConfiguration()
		if err != nil {
			fmt.Println("Aucune configuration trouvée")
			return
		}
		VMList, err := getVMlist(config)
		if err != nil {
			fmt.Println("Impossible de recupérer la liste des vm")
			fmt.Println(err)
			return
		}
		affichageVM(VMList)
	},
}

func init() {
	ServeurCmd.AddCommand(StartCmd)
}

func getVMlist(c *vmware.Client) (string, error) {
	req, err := http.NewRequest("GET", vmware.Host+"/rest/vcenter/vm", nil)
	if err != nil {
		return "", errors.New("erreur lors de la création de la requête GET")
	}
	req.Header.Add("vmware-api-session-id", c.Token)
	println("le token passé en authorization : ", c.Token)

	resp, err := vmware.InsecureHTTPClient.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de la requête GET :", err)
		return "", errors.New("Erreur lors de la requete http")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("erreur lors de la lecture du corps de la réponse")
	}
	if resp.StatusCode != 200 {
		return "truc", fmt.Errorf("requête échouée avec le code %d : %s", resp.StatusCode, string(body))
	}
	return string(body), nil
}

func affichageVM(jsonData string) error {
	var vmResp VMResponse
	err := json.Unmarshal([]byte(jsonData), &vmResp)
	if err != nil {
		return fmt.Errorf("erreur lors du parsing du JSON : %v", err)
	}

	if len(vmResp.Value) == 0 {
		fmt.Println("Aucune VM trouvée.")
		return nil
	}

	fmt.Println("Liste des VMs :")
	fmt.Println("ID        | Nom              | État")
	fmt.Println("--------------------------------------------")
	for _, vm := range vmResp.Value {
		fmt.Printf("%-10s | %-16s | %s\n", vm.VM, vm.Name, vm.PowerState)
	}
	return nil
}
