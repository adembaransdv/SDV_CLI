package vmware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"io"
	"gopkg.in/ini.v1"
)

var once sync.Once

const host string = "http://192.168.1.3:5000"

type client struct {
	Client_name string
	Token       string
}

type sessionResponse struct {
	Value string `json:"value"`
}

func DeleteVM(c *client, vmID string) error {
	url := host + "/rest/vcenter/vm/" + vmID
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.New("erreur lors de la création de la requête DELETE")
	}

	req.Header.Add("Authorization", c.Token)

	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de la requête DELETE : %v", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 204 {
		return fmt.Errorf("échec suppression VM, status %d : %s", resp.StatusCode, string(body))
	}

	return nil
}
func CheckConfiguration() (*client, error) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		return nil, fmt.Errorf("erreur lors de la lecture du fichier de configuration INI : %v", err)
	}

	login := cfg.Section("vmware").Key("login").String()
	password := cfg.Section("vmware").Key("password").String()

	if login == "" || password == "" {
		return nil, errors.New("les informations de connexion sont manquantes dans le fichier config.ini")
	}

	return Connection(login, password)
}

func Connection(login string, pass string) (*client, error) {
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("erreur lors de la lecture de la réponse")
	}

	var session sessionResponse
	err = json.Unmarshal(body, &session)
	if err != nil {
		return nil, errors.New("erreur lors du parsing JSON")
	}

	c := &client{
		Token: session.Value,
	}
	return c, nil
}
