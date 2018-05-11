package grpc

import (
	"context"
	"net/http"
	"sync"

	"google.golang.org/grpc"

	"golang.ysitd.cloud/log"

	api "code.ysitd.cloud/api/totp"

	"app.ysitd/auth/totp/pkg/service"
)

type Service struct {
	bootstrap sync.Once
	Logger    log.Logger       `inject:"grpc logger"`
	Server    *grpc.Server     `inject:""`
	Service   *service.Service `inject:""`
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.bootstrap.Do(func() {
		s.Server.GetServiceInfo()
		api.RegisterTotpServer(s.Server, s)
	})

	s.Server.ServeHTTP(w, r)
}

func (s *Service) IssueKey(ctx context.Context, in *api.IssueKeyRequest) (out *api.IssueKeyReply, err error) {
	url, recoverCode, err := s.Service.IssueKey(ctx, in.Issuer, in.Username)
	if err != nil {
		return
	}

	return &api.IssueKeyReply{
		Url:     url,
		Recover: recoverCode,
	}, nil
}

func (s *Service) Validate(ctx context.Context, in *api.ValidateRequest) (out *api.ValidateReply, err error) {
	validate, err := s.Service.ValidatePasscode(ctx, in.Issuer, in.Username, in.Passcode, *in.Time)
	if err != nil {
		return
	}

	return &api.ValidateReply{
		Validate: validate,
	}, nil
}

func (s *Service) RecoverKey(ctx context.Context, in *api.RecoverRequest) (out *api.RecoverReply, err error) {
	validate, url, recoverCode, err := s.Service.RecoverKey(ctx, in.Issuer, in.Username, in.Recover)
	if err != nil {
		return
	}

	return &api.RecoverReply{
		Validate: validate,
		Url:      url,
		Recover:  recoverCode,
	}, nil
}

func (s *Service) RemoveKey(ctx context.Context, in *api.RemoveKeyRequest) (out *api.RemoveKeyReply, err error) {
	err = s.Service.RemoveKey(ctx, in.Issuer, in.Username)
	if err != nil {
		return
	}

	return &api.RemoveKeyReply{
		Removed: true,
	}, nil
}
