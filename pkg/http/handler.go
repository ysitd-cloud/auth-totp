package http

import (
	"net/http"
	"strings"

	"app.ysitd/auth/totp/pkg/http/grpc"
	"app.ysitd/auth/totp/pkg/http/rest"
)

type Handler struct {
	RestHandler *rest.Server  `inject:""`
	GrpcHandler *grpc.Service `inject:""`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		h.GrpcHandler.ServeHTTP(w, r)
	} else {
		h.RestHandler.ServeHTTP(w, r)
	}
}
