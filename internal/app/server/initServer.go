package server

import (
	"github.com/gorilla/sessions"
	"net/http"
)

func Start(config *Config) error  {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s:= newServer(sessionStore)
	s.configureLogger(config.LogLevel)
	return http.ListenAndServe(config.BindAddr,s)
}