package grpc

import (
	"context"
	"net/http"
	"time"

	"code.ysitd.cloud/auth/totp/pkg/service"
	"golang.ysitd.cloud/log"
	"google.golang.org/grpc"

	api "code.ysitd.cloud/api/totp"
)

type Service struct {
	ready   bool
	Logger  log.Logger       `inject:"grpc logger"`
	Server  *grpc.Server     `inject:""`
	Service *service.Service `inject:""`
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !s.ready {
		s.Server.GetServiceInfo()
		api.RegisterTotpServer(s.Server, s)
		s.ready = true
	}

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
	timestamp := in.Time
	t := time.Unix(timestamp.Seconds, int64(timestamp.Nanos))
	validate, err := s.Service.ValidatePasscode(ctx, in.Issuer, in.Username, in.Passcode, t)
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
