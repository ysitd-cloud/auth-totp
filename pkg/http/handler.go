package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tonyhhyip/vodka"

	api "code.ysitd.cloud/api/totp"
)

func (s *Server) issueKey(c *vodka.Context) {
	defer c.Request.Body.Close()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		s.Logger.WithField("endpoint", "issueKey").Error(err)
		http.Error(c.Response, "error when reading body", http.StatusUnprocessableEntity)
		return
	}

	var req api.IssueKeyRequest

	if err := json.Unmarshal(body, &req); err != nil {
		s.Logger.WithField("endpoint", "issueKey").Error(err)
		http.Error(c.Response, "error when parsing body", http.StatusUnprocessableEntity)
		return
	}

	url, recoverCode, err := s.Service.IssueKey(c, req.Issuer, req.Username)
	if err != nil {
		s.Logger.WithField("endpoint", "issueKey").Error(err)
		http.Error(c.Response, "error when processing", 530)
		return
	}

	c.JSON(http.StatusCreated, api.IssueKeyReply{
		Url:     url,
		Recover: recoverCode,
	})
}

func (s *Server) validate(c *vodka.Context) {
	defer c.Request.Body.Close()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		s.Logger.WithField("endpoint", "validate").Error(err)
		http.Error(c.Response, "error when reading body", http.StatusUnprocessableEntity)
		return
	}

	var req struct {
		Issuer   string `json:"issuer"`
		Username string `json:"username"`
		Passcode string `json:"passcode"`
		Time     string `json:"time"`
	}

	if err := json.Unmarshal(body, &req); err != nil {
		s.Logger.WithField("endpoint", "validate").Error(err)
		http.Error(c.Response, "error when parsing body", http.StatusUnprocessableEntity)
		return
	}

	t, err := time.Parse(time.RFC3339, req.Time)

	if err := json.Unmarshal(body, &req); err != nil {
		s.Logger.WithField("endpoint", "validate").Error(err)
		http.Error(c.Response, "error when parsing time", http.StatusUnprocessableEntity)
		return
	}

	validate, err := s.Service.ValidatePasscode(c, req.Issuer, req.Username, req.Passcode, t.UTC())
	if err != nil {
		s.Logger.WithField("endpoint", "validate").Error(err)
		http.Error(c.Response, "error when processing", 530)
		return
	}

	if validate {
		c.Status(http.StatusOK)
	} else {
		c.Status(http.StatusUnauthorized)
	}
}

func (s *Server) recoverKey(c *vodka.Context) {
	defer c.Request.Body.Close()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		s.Logger.WithField("endpoint", "recoverKey").Error(err)
		http.Error(c.Response, "error when reading body", http.StatusUnprocessableEntity)
		return
	}

	var req api.RecoverRequest

	if err := json.Unmarshal(body, &req); err != nil {
		s.Logger.WithField("endpoint", "recoverKey").Error(err)
		http.Error(c.Response, "error when parsing body", http.StatusUnprocessableEntity)
		return
	}

	validate, url, recoverCode, err := s.Service.RecoverKey(c, req.Issuer, req.Username, req.Recover)

	if err != nil {
		s.Logger.WithField("endpoint", "recoverKey").Error(err)
		http.Error(c.Response, "error when processing", 530)
		return
	}

	if !validate {
		c.Status(http.StatusForbidden)
		return
	}

	c.JSON(http.StatusOK, api.IssueKeyReply{
		Url:     url,
		Recover: recoverCode,
	})
}

func (s *Server) removeKey(c *vodka.Context) {
	defer c.Request.Body.Close()

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		s.Logger.WithField("endpoint", "removeKey").Error(err)
		http.Error(c.Response, "error when reading body", http.StatusUnprocessableEntity)
		return
	}

	var req api.RemoveKeyRequest

	if err := json.Unmarshal(body, &req); err != nil {
		s.Logger.WithField("endpoint", "removeKey").Error(err)
		http.Error(c.Response, "error when parsing body", http.StatusUnprocessableEntity)
		return
	}

	err = s.Service.RemoveKey(c, req.Issuer, req.Username)

	if err != nil {
		s.Logger.WithField("endpoint", "removeKey").Error(err)
		http.Error(c.Response, "error when processing", 530)
		return
	}

	c.Status(http.StatusOK)
}
