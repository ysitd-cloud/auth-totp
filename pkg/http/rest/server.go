package rest

import (
	"net/http"

	"app.ysitd/auth/totp/pkg/service"
	"github.com/gorilla/handlers"
	"github.com/tonyhhyip/vodka"
	"golang.ysitd.cloud/log"
)

type Server struct {
	http.Handler
	Logger  log.Logger       `inject:"http logger"`
	Service *service.Service `inject:""`
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Handler == nil {
		s.initHandler()
	}

	s.Handler.ServeHTTP(w, r)
}

func (s *Server) initHandler() {
	router := vodka.NewRouter()

	router.POST("/key", s.issueKey)
	router.POST("/validate", s.validate)
	router.PUT("/key", s.recoverKey)
	router.DELETE("/key", s.removeKey)

	var handler http.Handler
	handler = vodka.CastHandlerForHTTP(router.Handler(), nil)
	handler = handlers.CombinedLoggingHandler(s.Logger.Writer(), handler)
	handler = handlers.RecoveryHandler(
		handlers.RecoveryLogger(s.Logger),
		handlers.PrintRecoveryStack(true),
	)(handler)
	s.Handler = handler
}
