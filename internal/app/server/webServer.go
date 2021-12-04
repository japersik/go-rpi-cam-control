package server

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"os"
	"time"
)

type server struct {
	router *mux.Router
	logger *zap.Logger
	atom   zap.AtomicLevel
	sessionStore sessions.Store
}

const (
	ctxKeyRequestId = iota
)

func newServer( sessionStore sessions.Store)  *server{
	s:= &server{
		router: mux.NewRouter(),
		atom:   zap.NewAtomicLevel(),
		sessionStore: sessionStore,
	}
	s.configureRouter()

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = ""
	s.logger = zap.New(zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), zapcore.Lock(os.Stdout), s.atom))

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
	s.router.ServeHTTP(w,r)
}

func (s *server) configureLogger(logLevel string) error {
	level := s.atom.Level()
	err := level.Set(logLevel)
	if err != nil {
		return err
	}
	return nil
}
func (s *server) configureRouter()  {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	s.router.PathPrefix("/static/").Handler( http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)

	//free sites
	s.router.HandleFunc("/", s.handleHomeHTML())
	s.router.HandleFunc("/register", s.handleRegisterHTML())
	s.router.HandleFunc("/createUser", s.handleCreateUser()).Methods("POST")
	s.router.HandleFunc("/authUser", s.handleMakeSession()).Methods("POST")
	s.router.HandleFunc("/logout", s.handleMakeSession()).Methods("POST")
	//sites with password
	private:= s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser) //аутентификация
	private.PathPrefix("/static/").Handler( http.StripPrefix("/private/static/", http.FileServer(http.Dir("./private/static/"))))
	private.HandleFunc("/", s.handleAlbumHTML())
	//private.HandleFunc("/control", s.handle())
}


func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		log:= s.logger.With(zap.String("remote_addr", request.RemoteAddr),
			zap.String("requestId", request.Context().Value(ctxKeyRequestId).(string)),
		)
		log.Info("started.",
			zap.String("Method",request.Method),
			zap.String("URI", request.RequestURI),
		)
		start:= time.Now()
		myWriter:= &responseWriter{writer,http.StatusOK}
		next.ServeHTTP(myWriter,request)
		log.Info("completed in",
			zap.Duration("Complete time",time.Now().Sub(start)),
			zap.Int("Status code:", myWriter.code),
			zap.String("Status:",http.StatusText(myWriter.code)),
		)
	})
}
func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("auth")
		next.ServeHTTP(writer,request)
	})
}
func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		id:= uuid.New().String()
		writer.Header().Set("X-Request-ID",id)
		next.ServeHTTP(writer,request.WithContext(context.WithValue(request.Context(),ctxKeyRequestId, id)))
	})
}

func (s *server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (s *server) handleMakeSession() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

