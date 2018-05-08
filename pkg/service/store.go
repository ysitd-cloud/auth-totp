package service

import (
	"context"
	"errors"
	"time"

	"github.com/facebookgo/inmem"
	"golang.ysitd.cloud/db"
)

var ErrUnexpectedRowAffected = errors.New("unexpected row affected")

const cacheExpire = time.Minute * 10

type Store struct {
	Opener *db.GeneralOpener `inject:""`
	Cache  inmem.Cache       `inject:"cache"`
}

func (s *Store) Create(ctx context.Context, issuer, username, secret string) (err error) {
	conn, err := s.Opener.Open()
	if err != nil {
		return
	}

	defer conn.Close()

	query := "INSERT INTO totp (issuer, username, secret) VALUES ($1, $2, $3)"

	result, err := conn.ExecContext(ctx, query, issuer, username, secret)
	if err != nil {
		return
	}

	if row, err := result.RowsAffected(); err != nil {
		return err
	} else if row != 1 {
		return ErrUnexpectedRowAffected
	}

	return nil
}

func cacheKey(issuer, username string) string {
	return username + "@" + issuer
}

func (s *Store) Get(ctx context.Context, issuer, username string) (secret string, err error) {
	key := cacheKey(issuer, username)
	if val, hit := s.Cache.Get(key); hit {
		return val.(string), nil
	}

	secret, err = s.getSecretFromDB(ctx, issuer, username)

	if err == nil {
		s.Cache.Add(key, secret, time.Now().Add(cacheExpire))
	}

	return secret, err
}

func (s *Store) getSecretFromDB(ctx context.Context, issuer, username string) (secret string, err error) {
	conn, err := s.Opener.Open()
	if err != nil {
		return
	}

	defer conn.Close()

	query := "SELECT secret FROM totp WHERE issuer = $1 AND username = $2"

	row := conn.QueryRowContext(ctx, query, issuer, username)

	err = row.Scan(&secret)

	return
}

func (s *Store) Update(ctx context.Context, issuer, username, secret string) (err error) {
	conn, err := s.Opener.Open()
	if err != nil {
		return
	}

	defer conn.Close()

	query := "UPDATE totp SET secret = $3 WHERE issuer = $1 AND username = $2,"

	result, err := conn.ExecContext(ctx, query, issuer, username, secret)
	if err != nil {
		return
	}

	if row, err := result.RowsAffected(); err != nil {
		return err
	} else if row != 1 {
		return ErrUnexpectedRowAffected
	}

	s.Cache.Remove(cacheKey(issuer, username))

	return nil
}

func (s *Store) Delete(ctx context.Context, issuer, username string) (err error) {
	conn, err := s.Opener.Open()
	if err != nil {
		return
	}

	defer conn.Close()

	query := "DELETE totp WHERE issuer = $1 AND username = $2,"

	result, err := conn.ExecContext(ctx, query, issuer, username)
	if err != nil {
		return
	}

	if row, err := result.RowsAffected(); err != nil {
		return err
	} else if row != 1 {
		return ErrUnexpectedRowAffected
	}

	s.Cache.Remove(cacheKey(issuer, username))

	return nil
}
