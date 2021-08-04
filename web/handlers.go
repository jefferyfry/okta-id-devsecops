package web

import (
	"github.com/jefferyfry/funclog"
	"net/http"
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







