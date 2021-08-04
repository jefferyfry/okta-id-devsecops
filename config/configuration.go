package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/jefferyfry/funclog"
	"os"
)

var (
	FrontendServiceEndpoint = "8086"
	HealthCheckEndpoint = "8096"
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	FrontendServiceEndpoint string `json:"frontendServiceEndpoint"`
	HealthCheckEndpoint string `json:"healthCheckEndpoint"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		FrontendServiceEndpoint,
		HealthCheckEndpoint,
	}

	if dir, err := os.Getwd(); err != nil {
		LogE.Println("Unable to determine working directory.")
		return conf, err
	} else {
		LogI.Printf("Running service with working directory %s \n", dir)
	}

	//parse commandline arguments
	configFile := flag.String("configFile", "", "set the path to the configuration json file")
	frontendServiceEndpoint := flag.String("frontendServiceEndpoint", "", "set the value of the frontend service endpoint port")
	healthCheckEndpoint := flag.String("healthCheckEndpoint", "", "set the value of the health check endpoint port")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CLOUD_BILL_FRONTEND_CONFIG_FILE")
	}
	if *frontendServiceEndpoint == "" {
		*frontendServiceEndpoint = os.Getenv("CLOUD_BILL_FRONTEND_SERVICE_ENDPOINT")
	}
	if *healthCheckEndpoint == "" {
		*healthCheckEndpoint = os.Getenv("CLOUD_BILL_FRONTEND_HEALTH_CHECK_ENDPOINT")
	}

	if *configFile == "" {
		//try other flags
		conf.FrontendServiceEndpoint = *frontendServiceEndpoint
		conf.HealthCheckEndpoint = *healthCheckEndpoint
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

	if conf.FrontendServiceEndpoint == "" {
		LogE.Println("FrontendServiceEndpoint was not set.")
		valid = false
	}

	if conf.HealthCheckEndpoint == "" {
		LogE.Println("HealthCheckEndpoint was not set.")
		valid = false
	}

	if gAppCredPath,gAppCredExists := os.LookupEnv("GOOGLE_APPLICATION_CREDENTIALS"); !gAppCredExists {
		LogE.Println("GOOGLE_APPLICATION_CREDENTIALS was not set. ")
		valid = false
	} else {
		if _, gAppCredPathErr := os.Stat(gAppCredPath); os.IsNotExist(gAppCredPathErr) {
			LogE.Println("GOOGLE_APPLICATION_CREDENTIALS file does not exist: ", gAppCredPath)
			valid = false
		} else {
			LogI.Println("Using GOOGLE_APPLICATION_CREDENTIALS file: ", gAppCredPath)
		}
	}

	if !valid {
		return conf, errors.New("Subscription frontend service configuration is not valid!")
	} else {
		return conf, nil
	}
}