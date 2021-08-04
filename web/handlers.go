package web

import (
	"github.com/jefferyfry/funclog"
	"github.com/okta/okta-jwt-verifier-golang"
	"net/http"
	"strings"
)

var (
	LogI = funclog.NewInfoLogger("INFO: ")
	LogE = funclog.NewErrorLogger("ERROR: ")
)

type ApiHandler struct {
}

func GetApiHandler() *ApiHandler {

	return &ApiHandler{}
}

func (hdlr *ApiHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (hdlr *ApiHandler) ValidateApiAccess(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - You are not authorized for this request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - You are authorized for this request"))
}



func isAuthenticated(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}
	tokenParts := strings.Split(authHeader, "Bearer ")
	bearerToken := tokenParts[1]

	toValidate := map[string]string{}
	toValidate["aud"] = "api://default"
	toValidate["cid"] = "{CLIENT_ID}"

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer: "https://${yourOktaDomain}/oauth2/default",
		ClaimsToValidate: toValidate,
	}

	_, err := jwtVerifierSetup.New().VerifyAccessToken(bearerToken)

	if err != nil {
		return false
	}

	return true
}







