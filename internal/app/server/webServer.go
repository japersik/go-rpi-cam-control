package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/japersik/go-rpi-cam-control/internal/app/service"
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
	service      service.Service
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

func newServer(config *Config, service service.Service) *server {
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := &server{
		router:       mux.NewRouter(),
		atom:         zap.NewAtomicLevel(),
		sessionStore: sessionStore,
		service:      service,
	}
	s.configureRouter()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""
	s.logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.Lock(os.Stdout), s.atom))
	s.configureLogger(config.LogLevel)
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
	private.HandleFunc("/{id:-?[0-9]+}", s.handleAlbumHTML()).Methods("GET")
	private.PathPrefix("/static/").Handler(http.StripPrefix("/private/static/", http.FileServer(http.Dir("./private/static/"))))
	private.HandleFunc("/move_control", s.moveControl()).Methods("POST")
	private.HandleFunc("/camera_control", s.cameraControl()).Methods("POST")
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
