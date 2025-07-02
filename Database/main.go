package database

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

	"gopkg.in/ini.v1"
)

func CheckDatabase() (bool, error) {
	filename := "Database.db"

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		cfg := ini.Empty()
		if err := cfg.SaveTo(filename); err != nil {
			return false, fmt.Errorf("erreur création fichier : %v", err)
		}
		fmt.Println("Fichier créé :", filename)
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("impossible de lire le fichier : %v", err)
	}

	_, err = base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		encoded := base64.StdEncoding.EncodeToString(data)
		err = os.WriteFile(filename, []byte(encoded), 0644)
		if err != nil {
			return false, fmt.Errorf("erreur écriture fichier encodé : %v", err)
		}
		fmt.Println("Fichier Database.db encodé en Base64")
	} else {
		fmt.Println("Fichier déjà encodé en Base64")
	}

	return true, nil
}

func FindInBDD(key string) (bool, error) {
	filename := "Database.db"

	encodedData, err := os.ReadFile(filename)
	if err != nil {
		return false, fmt.Errorf("erreur lecture fichier base64 : %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		return false, fmt.Errorf("erreur décodage base64 : %v", err)
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return false, fmt.Errorf("erreur chargement ini : %v", err)
	}

	val := cfg.Section("default").Key(key).String()
	if val == "" {
		return false, nil
	}

	return true, nil
}

func AddKeyToBDD(key, value string) error {
	filename := "Database.db"
	encodedData, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("erreur lecture fichier base64 : %v", err)
	}

	data, err := base64.StdEncoding.DecodeString(string(encodedData))
	if err != nil {
		return fmt.Errorf("erreur décodage base64 : %v", err)
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return fmt.Errorf("erreur chargement ini : %v", err)
	}

	cfg.Section("default").Key(key).SetValue(value)

	var buf bytes.Buffer
	_, err = cfg.WriteTo(&buf)
	if err != nil {
		return fmt.Errorf("erreur écriture ini en mémoire : %v", err)
	}

	encoded := base64.StdEncoding.EncodeToString(buf.Bytes())

	err = os.WriteFile(filename, []byte(encoded), 0644)
	if err != nil {
		return fmt.Errorf("erreur écriture fichier encodé : %v", err)
	}

	return nil
}
