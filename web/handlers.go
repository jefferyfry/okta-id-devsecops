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
	Aud string
	Cid string
	Domain string
}

func GetApiHandler(aud string,cid string,domain string) *ApiHandler {
	return &ApiHandler{aud,
		cid,
	domain}
}

func (hdlr *ApiHandler) Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (hdlr *ApiHandler) ValidateApiAccess(w http.ResponseWriter, r *http.Request) {
	if !isAuthenticated(r,hdlr.Aud,hdlr.Cid,hdlr.Domain) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - You are not authorized for this request"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 - You are authorized for this request"))
}



func isAuthenticated(r *http.Request,aud string, cid string, domain string) bool {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return false
	}
	tokenParts := strings.Split(authHeader, "Bearer ")
	bearerToken := tokenParts[1]

	toValidate := map[string]string{}
	toValidate["aud"] = aud
	toValidate["cid"] = cid

	jwtVerifierSetup := jwtverifier.JwtVerifier{
		Issuer: "https://"+domain+"/oauth2/default",
		ClaimsToValidate: toValidate,
	}

	_, err := jwtVerifierSetup.New().VerifyAccessToken(bearerToken)

	if err != nil {
		return false
	}

	return true
}







