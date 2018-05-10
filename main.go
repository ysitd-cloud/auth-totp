package main

import (
	"context"
	"net/http"
	"time"

	"code.ysitd.cloud/auth/totp/pkg/bootstrap"
	"os"
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
