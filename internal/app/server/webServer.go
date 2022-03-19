package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/japersik/go-rpi-cam-control/internal/app/cameraController"
	"github.com/japersik/go-rpi-cam-control/internal/app/moveController"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

type server struct {
	router       *mux.Router
	logger       *zap.Logger
	atom         zap.AtomicLevel
	sessionStore sessions.Store
	mover        moveController.Mover
	camera       cameraController.Camera
	userStore    userStore
}

const (
	ctxKeyRequestId = iota
	ctxKeyUser
)

var (
	errEmailOrPassword  = errors.New("Incorrect Email or password Error")
	errNotAuthenticated = errors.New("Not Authenticated")
	sessionName         = "SessionCookie"
)

func newServer(config *Config, mover moveController.Mover, camera cameraController.Camera) *server {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := &server{
		router:       mux.NewRouter(),
		atom:         zap.NewAtomicLevel(),
		sessionStore: sessionStore,
		mover:        mover,
		camera:       camera,
		userStore:    newUserStore(config.AuthFile),
	}
	s.configureRouter()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""
	s.logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.Lock(os.Stdout), s.atom))
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureLogger(logLevel string) error {
	level := s.atom.Level()
	err := level.Set(logLevel)
	if err != nil {
		return err
	}
	return nil
}

// configureRouter настройка обработчиков запросов
func (s *server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	//free sites
	s.router.HandleFunc("/", s.handleHomeHTML())
	s.router.HandleFunc("/register", s.handleRegisterHTML()).Methods("GET")
	s.router.HandleFunc("/register", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/login", s.handleLoginHTML()).Methods("GET")
	s.router.HandleFunc("/login", s.handleMakeSession()).Methods("POST")
	s.router.HandleFunc("/logout", s.logoutUser()).Methods("POST")
	//sites with password
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser) //аутентификация
	private.HandleFunc("/", s.handleAlbumHTML()).Methods("GET")
	private.PathPrefix("/static/").Handler(http.StripPrefix("/private/static/", http.FileServer(http.Dir("./private/static/"))))
	s.router.HandleFunc("/move_control", s.moveControl()).Methods("POST")
	s.router.HandleFunc("/camera_control", s.cameraControl()).Methods("POST")
	//private.HandleFunc("/control", s.handle())
}

//setRequestID ..
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id := uuid.New().String()
		writer.Header().Set("X-Request-ID", id)
		next.ServeHTTP(writer, request.WithContext(context.WithValue(request.Context(), ctxKeyRequestId, id)))
	})
}

//logRequest логирование запросов
func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		log := s.logger.With(zap.String("remote_addr", request.RemoteAddr),
			zap.String("requestId", request.Context().Value(ctxKeyRequestId).(string)),
		)
		log.Info("started",
			zap.String("Method", request.Method),
			zap.String("URI", request.RequestURI),
		)
		start := time.Now()
		myWriter := &responseWriter{writer, http.StatusOK}
		next.ServeHTTP(myWriter, request)
		log.Info("completed",
			zap.Duration("Complete time", time.Now().Sub(start)),
			zap.Int("Status code:", myWriter.code),
			zap.String("Status:", http.StatusText(myWriter.code)),
		)
	})
}

//authenticateUser аутентификация пользователя
func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			//fmt.Println(err)
			s.error(w, r, http.StatusInternalServerError, err)
			//http.Redirect(w, r, "/", http.StatusPermanentRedirect)
			return
		}
		uname, ok := session.Values["user_name"]
		if !ok {
			//s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		if !s.userStore.checkUserExist(uname.(string)) {
			//s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, uname)))
		//next.ServeHTTP(w, r)
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
		if req.UserCreationCodeWord != "passcode" || s.userStore.checkUserExist(req.Name) {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("user exist or wrong passcode"))
			return
		}
		err := s.userStore.CreateUser(req.Name, req.Password)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("user exist or wrong passcode"))
			return
		}
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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

		if !s.userStore.CheckUser(req.Name, req.Password) {
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
		fmt.Println(req)
		switch req.CommandName {
		case "take_photo":
			s.camera.TakePhoto()
		case "del_photo":
			s.camera.DelPhoto(req.id)
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
			s.mover.MoveX(req.Value)
		case "y":
			s.mover.MoveY(req.Value)
		}
		s.respond(w, r, http.StatusOK, "Ok")
	}
}

//error возвращает в http.Request ошибку с заданным кодом
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

//respond возвращает в http.Request в ответ в формате json
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
