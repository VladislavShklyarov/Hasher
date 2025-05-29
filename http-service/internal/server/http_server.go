package server

import (
	"fmt"
	"http-service/internal/config"

	"http-service/internal/app"
	"http-service/internal/transport/http"
	stdHttp "net/http"
	"time"
)

func RunHttpServer(app *app.Clients, cfg *config.Config) {

	router := http.NewRouter(app)
	server := &stdHttp.Server{
		Addr:           cfg.HttpAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("server started on: %s", cfg.HttpAddr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
