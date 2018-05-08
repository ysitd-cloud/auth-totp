package bootstrap

import (
	"net/http"

	"code.ysitd.cloud/auth/totp/pkg/grpc"
	httpService "code.ysitd.cloud/auth/totp/pkg/http"
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
)

var httpSrv httpService.Server
var grpcSrv grpc.Service

func init() {
	var graph inject.Graph
	graph.Logger = initLogger()

	graph.Provide(
		&inject.Object{Value: &httpSrv},
		&inject.Object{Value: &grpcSrv},
	)

	for _, fn := range []func(*inject.Graph){
		injectLogger,
		injectGrpc,
		injectStore,
	} {
		fn(&graph)
	}

	if err := graph.Populate(); err != nil {
		logger.Error(err)
		panic(err)
	}
}

func GetHttpServer() http.Handler {
	return &httpSrv
}

func GetGrpcServer() http.Handler {
	return &grpcSrv
}

func GetMainLogger() logrus.FieldLogger {
	return logger.WithField("source", "main")
}
