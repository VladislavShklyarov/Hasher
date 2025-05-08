package main

import (
	"fmt"
	"http-service/internal/transport/http"
	stdHttp "net/http"
	"time"
)

func main() {
	router := http.NewRouter()
	server := &stdHttp.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    120 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("server started on :8080")
}
