package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"app.ysitd/auth/totp/pkg/bootstrap"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "50051"
	}

	logger := bootstrap.GetMainLogger()
	server := http.Server{Handler: bootstrap.GetMainHandler(), Addr: ":" + port}

	logger.Debugf("Start HTTP Listen at %s", port)
	if err := server.ListenAndServe(); err != nil {
		logger.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	server.Shutdown(ctx)
}
