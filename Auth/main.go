package vmware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"gopkg.in/ini.v1"
)

var once sync.Once

type client struct {
	Client_name string
	Token       string
}

type sessionResponse struct {
	Value string `json:"value"`
}

type VM struct {
	VM         string `json:"vm"`
	Name       string `json:"name"`
	PowerState string `json:"power_state"`
}

type VMResponse struct {
	Value []VM `json:"value"`
}

const host string = "http://127.0.0.1:5000"

func Connection(login string, pass string) (*client, error) {
	// APPEL API POUR GET LE TOKEN
	url := host + "/rest/com/vmware/cis/session"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, errors.New("erreur lors de la création de la requête")
	}

	auth := login + ":" + pass
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+authEncoded)

	clienthttp := &http.Client{}
	resp, err := clienthttp.Do(req)
	if err != nil {
		fmt.Println("Erreur lors de la requête POST :", err)
		return nil, errors.New("erreur lors de l'appel HTTP")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("erreur lors de la lecture de la réponse")
	}

	var session sessionResponse
	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, errors.New("erreur lors du parsing JSON")
	}

	c := &client{
		Client_name: login,
		Token:       session.Value,
	}
	return c, nil
}

func GetVMlist(c *client) (string, error) {
	req, err := http.NewRequest("GET", host+"/rest/vcenter/vm", nil)
	if err != nil {
		return "", errors.New("erreur lors de la création de la requête GET")
	}
	req.Header.Add("Authorization", c.Token)

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
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

func AffichageVM(jsonData string) error {
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

func CheckConfiguration() (*client, error) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Erreur de lecture du fichier de configuration INI : %v", err)
	}

	login := cfg.Section("vmware").Key("login").String()
	password := cfg.Section("vmware").Key("password").String()

	return Connection(login, password)

}
