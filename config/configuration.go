package config

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/jefferyfry/funclog"
	"os"
)

var (
	ServiceEndpoint = "8086"
	HealthCheckEndpoint = "8096"
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ServiceConfig struct {
	ServiceEndpoint string `json:"serviceEndpoint"`
	HealthCheckEndpoint string `json:"healthCheckEndpoint"`
}

func GetConfiguration() (ServiceConfig, error) {
	conf := ServiceConfig {
		ServiceEndpoint,
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
	serviceEndpoint := flag.String("serviceEndpoint", "", "set the value of the service endpoint port")
	healthCheckEndpoint := flag.String("healthCheckEndpoint", "", "set the value of the health check endpoint port")
	flag.Parse()

	//try environment variables if necessary
	if *configFile == "" {
		*configFile = os.Getenv("CONFIG_FILE")
	}
	if *serviceEndpoint == "" {
		*serviceEndpoint = os.Getenv("SERVICE_ENDPOINT")
	}
	if *healthCheckEndpoint == "" {
		*healthCheckEndpoint = os.Getenv("HEALTH_CHECK_ENDPOINT")
	}

	if *configFile == "" {
		//try other flags
		conf.ServiceEndpoint = *serviceEndpoint
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

	if conf.ServiceEndpoint == "" {
		LogE.Println("ServiceEndpoint was not set.")
		valid = false
	}

	if conf.HealthCheckEndpoint == "" {
		LogE.Println("HealthCheckEndpoint was not set.")
		valid = false
	}

	if !valid {
		return conf, errors.New("Subscription frontend service configuration is not valid!")
	} else {
		return conf, nil
	}
}