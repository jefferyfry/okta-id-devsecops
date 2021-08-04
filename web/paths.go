package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

//SetUpService sets up the subscription service.
func SetUpService(webServiceEndpoint string,healthCheckEndpoint string) error {
	handler := GetApiHandler()

	healthCheck := mux.NewRouter()
	healthCheck.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)
	go http.ListenAndServe(":"+healthCheckEndpoint, healthCheck)

	webService := mux.NewRouter()

	webService.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)

	return http.ListenAndServe(":"+webServiceEndpoint, webService)
}