package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"code.ysitd.cloud/auth/totp/pkg/bootstrap"
)

func main() {
	logger := bootstrap.GetMainLogger()
	httpServer := http.Server{
		Addr:    ":50050",
		Handler: bootstrap.GetHttpServer(),
	}

	grpcServer := http.Server{
		Addr:    ":50051",
		Handler: bootstrap.GetGrpcServer(),
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		logger.Debugln("Start HTTP Listen")
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()

	go func() {
		logger.Debugln("Start Grpc Listen")
		if err := grpcServer.ListenAndServe(); err != nil {
			logger.Error(err)
		}
	}()

	<-c
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	httpServer.Shutdown(ctx)
	grpcServer.Shutdown(ctx)
}
