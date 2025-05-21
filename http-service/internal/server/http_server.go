package server

import (
	"fmt"

	"http-service/internal/app"
	"http-service/internal/transport/http"
	stdHttp "net/http"
	"time"
)

func RunHttpServer(app *app.Clients) {

	router := http.NewRouter(app)
	server := &stdHttp.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("server started on :8080")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
