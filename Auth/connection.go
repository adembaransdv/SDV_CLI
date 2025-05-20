
package vmware

import (
    "encoding/base64"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "net/http"
)

type Client struct {
    Username string
    Token    string
}

type sessionResponse struct {
    Value string `json:"value"`
}

func Connection(login string, pass string) (*Client, error) {
    url := "https://192.168.1.3/rest/com/vmware/cis/session"

    req, err := http.NewRequest("POST", url, nil)
    if err != nil {
        return nil, errors.New("erreur lors de la création de la requête")
    }

    encoded := base64.StdEncoding.EncodeToString([]byte(login + ":" + pass))
    req.Header.Add("Authorization", "Basic " + encoded)
    req.Header.Set("Content-Type", "application/json")

    clientHTTP := &http.Client{}
    resp, err := clientHTTP.Do(req)
    if err != nil {
        return nil, errors.New("erreur lors de la requête POST")
    }
    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var session sessionResponse
    if err := json.Unmarshal(body, &session); err != nil {
        return nil, errors.New("erreur parsing JSON")
    }

    if session.Value == "" {
        return nil, errors.New("token non reçu")
    }

    return &Client{
        Username: login,
        Token:    session.Value,
    }, nil
}
