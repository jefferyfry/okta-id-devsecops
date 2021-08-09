package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/jefferyfry/funclog"
	"os"
)

var (
	ServiceEndpoint = "80"
	Aud = "api.acme.com/test"
	Cid = "0oa1emw7xmqeh4Spd5d7"
	Issuer = "https://dev-73225252.okta.com/oauth2/aus1efvp3jwospP0Y5d7"
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	ServiceEndpoint string `json:"serviceEndpoint"`
	Aud string `json:"aud"`
	Cid string `json:"cid"`
	Issuer string `json:"issuer"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig{
		ServiceEndpoint,
		Aud,
		Cid,
		Issuer,
	}

	if dir, err := os.Getwd(); err != nil {
		LogE.Println("Unable to determine working directory.")
		return conf, err
	} else {
		LogI.Printf("Running service with working directory %s \n", dir)
	}

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	serviceEndpoint := flag.String("serviceEndpoint", "", "set the value of the service endpoint port")
	aud := flag.String("aud", "", "set the value of the audience")
	cid := flag.String("cid", "", "set the value of the client id")
	issuer := flag.String("issuer", "", "set the value of the auth server domain")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CONFIG_FILE")
	}
	if *serviceEndpoint == "" {
		*serviceEndpoint = os.Getenv("SERVICE_ENDPOINT")
	}
	if *aud == "" {
		*aud = os.Getenv("AUD")
	}
	if *cid == "" {
		*cid = os.Getenv("CID")
	}
	if *issuer == "" {
		*issuer = os.Getenv("ISSUER")
	}

	if *configFile == "" {
		//try other flags
		conf.ServiceEndpoint = *serviceEndpoint
		conf.Aud = *aud
		conf.Cid = *cid
		conf.Issuer = *issuer
	} else {
		if file, err := os.Open(*configFile); err != nil {
			LogE.Printf("Error reading confile file %s %s", *configFile, err)
			return conf, err
		} else {
			if err = json.NewDecoder(file).Decode(&conf); err != nil {
				return conf, errors.New("Configuration file not found.")
			}
			LogI.Printf("Using confile file %s \n", *configFile)
		}
	}

	valid := true

	if conf.ServiceEndpoint == "" {
		LogE.Println("ServiceEndpoint was not set.")
		valid = false
	}

	if conf.Aud == "" {
		LogE.Println("Aud was not set.")
		valid = false
	}

	if conf.Cid == "" {
		LogE.Println("Cid was not set.")
		valid = false
	}

	if conf.Issuer == "" {
		LogE.Println("Issuer was not set.")
		valid = false
	}

	if !valid {
		return conf, errors.New("api service configuration is not valid")
	} else {
		return conf, nil
	}
}