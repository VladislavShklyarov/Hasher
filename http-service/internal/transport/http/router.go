package http

import "github.com/julienschmidt/httprouter"

func NewRouter() *httprouter.Router {
	router := httprouter.New()

	router.POST("/process", ProcessDataHandler)

	return router
}
