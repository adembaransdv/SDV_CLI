package vmware

import (
	"crypto/tls"
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

type Client struct {
	Client_name string
	Token       string
}

type sessionResponse struct {
	Value string `json:"value"`
}

const Host string = "https://192.168.1.3"

var InsecureHTTPClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func Connection(login string, pass string) (*Client, error) {
	// APPEL API POUR GET LE TOKEN
	url := Host + "/rest/com/vmware/cis/session"

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, errors.New("erreur lors de la création de la requête")
	}

	auth := login + ":" + pass
	authEncoded := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Add("Authorization", "Basic "+authEncoded)

	// utilise le client global
	resp, err := InsecureHTTPClient.Do(req)
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

	return &Client{
		Client_name: login,
		Token:       session.Value,
	}, nil
}

func CheckConfiguration() (*Client, error) {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Erreur de lecture du fichier de configuration INI : %v", err)
	}

	login := cfg.Section("vmware").Key("login").String()
	password := cfg.Section("vmware").Key("password").String()

	return Connection(login, password)

}
