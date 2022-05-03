package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Secrets struct {
	DBConnectionString string `json:"db_conn_string"`
	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
}

func LoadSecrets() Secrets {

	jsonFile, err := os.Open("secrets.json")

	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	secrets := Secrets{}

	err = json.Unmarshal(byteValue, &secrets)
	if err != nil {
		panic(err)
	}

	return secrets
}
