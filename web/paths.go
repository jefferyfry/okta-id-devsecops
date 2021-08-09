package web

import (
	"github.com/gorilla/mux"
	"net/http"
)

//SetUpService sets up the API service.
func SetUpService(webServiceEndpoint string,aud string, cid string, domain string) error {
	handler := GetApiHandler(aud,cid,domain)

	apiService := mux.NewRouter()

	apiService.Methods(http.MethodGet).Path("/healthz").HandlerFunc(handler.Healthz)

	apiService.Methods(http.MethodGet).Path("/api/v1/users").HandlerFunc(handler.ValidateApiAccess)

	return http.ListenAndServe(":"+webServiceEndpoint, apiService)
}