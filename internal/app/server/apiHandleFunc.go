package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

//authenticateUser аутентификация пользователя
func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		uname, ok := session.Values["user_name"]
		if !ok {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		if !s.service.CheckUserExist(uname.(string)) {
			session.Options.MaxAge = -1
			err = s.sessionStore.Save(r, w, session)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			} else {
				s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			}
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, uname)))
	})
}

//handleCreateUser создаёт пользователя и в случае успеха переадресовывает на авторизацию
func (s *server) handleCreateUser() http.HandlerFunc {
	type request struct {
		UserCreationCodeWord string `json:"code_word"`
		Name                 string `json:"name"`
		Password             string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if s.service.CheckUserExist(req.Name) {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("user exist or wrong passcode"))
			return
		}
		err := s.service.CreateUser(req.Name, req.Password, req.UserCreationCodeWord)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("user exist or wrong passcode"))
			return
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
}

//handleMakeSession авторизовывает пользователя и устанавливает сессию
func (s *server) handleMakeSession() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if !s.service.CheckUser(req.Name, req.Password) {
			s.error(w, r, http.StatusUnauthorized, errEmailOrPassword)
			return
		} else {
			session, err := s.sessionStore.Get(r, sessionName)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
			session.Values["user_name"] = req.Name
			session.Options.MaxAge = int(time.Hour.Seconds())
			err = s.sessionStore.Save(r, w, session)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		//http.Redirect(w, r, "/", http.StatusPermanentRedirect)
		s.respond(w, r, http.StatusOK, "That Fits")
	}
}

//logoutUser сброс авторизации пользователя
func (s *server) logoutUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		_, ok := session.Values["user_name"]
		if ok {
			session.Options.MaxAge = -1
			err = s.sessionStore.Save(r, w, session)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
			}
		}
		s.respond(w, r, http.StatusNonAuthoritativeInfo, "Log out")
	}
}

func (s *server) cameraControl() http.HandlerFunc {
	type request struct {
		CommandName string `json:"command_name"`
		id          int    `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		switch req.CommandName {
		case "take_photo":
			ph, _ := s.service.TakePhoto()
			s.respond(w, r, http.StatusOK, ph)
			return
		case "del_photo":
			s.service.DelPhoto(req.id)
		case "get_photo":
			ph, _ := s.service.GetPhoto(req.id)
			s.respond(w, r, http.StatusOK, ph)
			return
		}

		s.respond(w, r, http.StatusOK, "Ok")
	}
}

//moveControl
func (s *server) moveControl() http.HandlerFunc {
	type request struct {
		Direction string `json:"direction"`
		Value     int    `json:"value"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		fmt.Println(req)
		switch req.Direction {
		case "x":
			s.service.MoveX(req.Value)
		case "y":
			s.service.MoveY(req.Value)
		}
		s.respond(w, r, http.StatusOK, "Ok")
	}
}
