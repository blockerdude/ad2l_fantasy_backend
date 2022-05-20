package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Secrets Secrets `json:"secrets"`
	OIDC    OIDC    `json:"oidc"`
}

type Secrets struct {
	DBConnectionString string `json:"db_conn_string"`
	GoogleClientID     string `json:"google_client_id"`
	GoogleClientSecret string `json:"google_client_secret"`
}

type OIDC struct {
	LoginRedirectURL  string `json:"login_redirect_url"`
	SignupRedirectURL string `json:"signup_redirect_url"`
	UIBaseURL         string `json:"ui_base_url"`
}

func LoadSecrets() Config {

	jsonFile, err := os.Open("conf.json")

	if err != nil {
		panic(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	config := Config{}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}

	return config
}
