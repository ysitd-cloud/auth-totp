package bootstrap

import (
	"net/http"

	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"

	httpService "app.ysitd/auth/totp/pkg/http"
)

var handler httpService.Handler

func init() {
	var graph inject.Graph
	graph.Logger = initLogger().WithField("source", "inject")

	graph.Provide(
		&inject.Object{Value: &handler},
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

func GetMainHandler() http.Handler {
	return &handler
}

func GetMainLogger() logrus.FieldLogger {
	return logger.WithField("source", "main")
}
