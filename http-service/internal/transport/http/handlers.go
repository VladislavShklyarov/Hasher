package http

import (
	"github.com/julienschmidt/httprouter"
	"http-service/internal/service"
	"net/http"
)

func ProcessDataHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	body := service.ValidateHttpRequest(w, r)
	if len(body) != 0 {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Recieved: " + string(body)))
	}
}
