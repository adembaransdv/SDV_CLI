package vmware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io" 
	"net/http"
	"gopkg.in/ini.v1"
)

const host string = "http://127.0.0.1:5000"

type client struct {
	Token string
}

type sessionResponse struct {
	Value string `json:"value"`
}

func StopVM(c *client, vmID string) error {
	url := fmt.Sprintf("%s/rest/vcenter/vm/%s/power/stop", host, vmID)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return errors.New("erreur lors de la création de la requête POST")
	}
	req.Header.Add("Authorization", "Bearer "+c.Token)
	clientHTTP := &http.Client{}
	resp, err := clientHTTP.Do(req)
	if err != nil {
		return fmt.Errorf("erreur lors de l'appel HTTP : %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("erreur lors de l'arrêt de la VM, code : %d", resp.StatusCode)
	}
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("erreur lors du décodage de la réponse : %v", err)
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
