package api

import (
	"net/http"
	"time"

	"github.com/EmptyLungs/golang_file_server/pkg/files"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Config struct {
	Host                  string        `mapstructure:"host"`
	Port                  string        `mapstructure:"port"`
	HttpServerTimeout     time.Duration `mapstructure:"http-server-timeout"`
	UploaderDir           string        `mapstructure:"upload-dir"`
	UploaderMaxFileSizeMB int64         `mapstructure:"upload-max-file-size"`
	AuthToken             string        `mapstructure:"auth-token"`
}

type Server struct {
	router      *mux.Router
	logger      *zap.Logger
	config      *Config
	handler     http.Handler
	fileManager files.IFileManager
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/echo", s.echoHandler).Methods("GET")
	s.router.HandleFunc("/upload", s.uploadFileHandler).Methods("POST")
	s.router.HandleFunc("/list", s.listFileHandler).Methods("GET")
	s.router.HandleFunc("/delete", s.deleteFileHandler).Methods("POST")
}

func (s *Server) registerMiddlewares() {
	httpLogger := NewLoggingMiddleware(s.logger)
	authMiddleware := NewAuthMiddleware(s.config.AuthToken)
	s.router.Use(httpLogger.Handler)
	s.router.Use(authMiddleware.Handler)
}

func (s *Server) ListenAndServe() {
	// todo: add ready state
	s.startServer()
}

func (s *Server) startServer() {
	srv := &http.Server{
		Addr:         s.config.Host + ":" + s.config.Port,
		WriteTimeout: s.config.HttpServerTimeout,
		ReadTimeout:  s.config.HttpServerTimeout,
		IdleTimeout:  2 * s.config.HttpServerTimeout,
		Handler:      s.handler,
	}

	s.logger.Info("start_server", zap.String("addr", srv.Addr))
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		s.logger.Fatal("HTTP server crashed", zap.Error(err))
	}
}

func NewServer(config *Config, logger *zap.Logger, fileManager files.IFileManager) (*Server, error) {
	childLogger := logger.With(zap.String("source", "server"))
	srv := &Server{
		router:      mux.NewRouter(),
		logger:      childLogger,
		config:      config,
		fileManager: fileManager,
	}
	srv.registerHandlers()
	srv.registerMiddlewares()
	srv.handler = srv.router

	return srv, nil
}
