package main

import (
	"github.com/jefferyfry/funclog"
	"okta-id-devsecops/config"
	"okta-id-devsecops/web"
)

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

func main() {
	LogI.Println("Starting API Service...")
	conf, err := config.GetConfiguration()

	if err != nil {
		LogE.Fatalf("Invalid configuration: %v", err)
	}

	//start web service
	LogE.Fatal(web.SetUpService(conf.ServiceEndpoint,conf.HealthCheckEndpoint,conf.Aud,conf.Cid, config.Domain))
}