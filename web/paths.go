package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

//SetUpService sets up the API service.
func SetUpService(webServiceEndpoint string,healthCheckEndpoint string) error {
	handler := GetApiHandler()

	apiService := mux.NewRouter()

	apiService.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)

	apiService.Methods(http.MethodGet).Path("/api/users").HandlerFunc(handler.ValidateApiAccess)

	return http.ListenAndServe(":"+webServiceEndpoint, apiService)
}