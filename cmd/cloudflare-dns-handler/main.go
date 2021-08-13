package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/letsencrypt/sre-tools/cloudflare-dns-handler/cloudflare"
	"github.com/letsencrypt/sre-tools/cloudflare-dns-handler/config"
)

func main() {
	token, err := getAPIToken()
	if err != nil {
		log.Fatalln(err)
	}

	// TODO(@aaomidi): nice error messages if arg is not provided
	configPath := os.Args[1]

	conf, err := config.ReadConfigurationFile(configPath)
	if err != nil {
		log.Fatalln(err)
	}

	cf := cloudflare.Init(token)

	err = cf.Apply(conf)
	if err != nil {
		log.Fatalln(err)
	}
}

func getAPIToken() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	path := path.Join(homeDir, ".config", "cloudflare-dns-handler", "api_token")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
