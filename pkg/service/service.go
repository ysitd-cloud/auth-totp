package service

import (
	"context"
	"encoding/base32"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
	"golang.ysitd.cloud/log"
)

var validateOption = totp.ValidateOpts{
	Period:    30,
	Skew:      1,
	Digits:    otp.DigitsSix,
	Algorithm: otp.AlgorithmSHA256,
}

type Service struct {
	Store  *Store     `inject:""`
	Logger log.Logger `inject:"service logger"`
}

func (s *Service) IssueKey(ctx context.Context, issuer, username string) (url string, recover string, err error) {
	return s.issueKey(ctx, issuer, username, true)
}

func (s *Service) issueKey(ctx context.Context, issuer, username string, newIssue bool) (url string, recover string, err error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      issuer,
		AccountName: username,
		Algorithm:   otp.AlgorithmSHA256,
	})

	if err != nil {
		return
	}

	secret := key.Secret()

	if newIssue {
		err = s.Store.Create(ctx, issuer, username, secret)
	} else {
		err = s.Store.Update(ctx, issuer, username, secret)
	}
	if err != nil {
		return
	}

	recoverCode, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	recover = base32.HexEncoding.EncodeToString(recoverCode)

	return key.String(), recover, nil
}

func (s *Service) ValidatePasscode(ctx context.Context, issuer, username, passcode string, t time.Time) (validate bool, err error) {
	s.Logger.Debugf("Validate passcode %s:%s@%s", username, passcode, issuer)
	secret, err := s.Store.Get(ctx, issuer, username)
	if err != nil {
		return
	}

	return totp.ValidateCustom(passcode, secret, t, validateOption)
}

func (s *Service) RecoverKey(ctx context.Context, issuer, username, recover string) (validate bool, url, newRecover string, err error) {
	secret, err := s.Store.Get(ctx, issuer, username)
	if err != nil {
		return
	}

	var hash []byte
	_, err = base32.HexEncoding.Decode(hash, []byte(recover))
	if err != nil {
		return
	}

	validate = bcrypt.CompareHashAndPassword(hash, []byte(secret)) != nil

	url, newRecover, err = s.issueKey(ctx, issuer, username, false)

	return
}

func (s *Service) RemoveKey(ctx context.Context, issuer, username string) error {
	return s.Store.Delete(ctx, issuer, username)
}
