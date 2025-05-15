package http

import (
	"github.com/julienschmidt/httprouter"
	"github.com/swaggo/http-swagger"
	_ "http-service/cmd/docs"
	"http-service/internal/app"
	handlers "http-service/internal/transport/http/handlers"
	"net/http"
)

func NewRouter(app *app.Clients) *httprouter.Router {
	router := httprouter.New()

	router.POST("/process", handlers.ProcessDataHandler(app))
	router.GET("/getLog", handlers.ReadLogHandler(app))
	router.DELETE("/deleteLog", handlers.DeleteLogHandler(app))
	router.GET("/swagger/*any", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		httpSwagger.WrapHandler.ServeHTTP(w, r)
	})

	return router
}
