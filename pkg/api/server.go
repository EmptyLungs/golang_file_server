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
}

type Server struct {
	router      *mux.Router
	logger      *zap.Logger
	config      *Config
	handler     http.Handler
	fileManager files.FileManager
}

func (s *Server) registerHandlers() {
	s.router.HandleFunc("/echo", s.echoHandler).Methods("GET")
	s.router.HandleFunc("/upload", s.uploadFileHandler).Methods("POST")
	s.router.HandleFunc("/list", s.listFileHandler).Methods("GET")
}

func (s *Server) registerMiddlewares() {
	httpLogger := NewLoggingMiddleware(s.logger)
	s.router.Use(httpLogger.Handler)
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

func NewServer(config *Config, logger *zap.Logger) (*Server, error) {
	childLogger := logger.With(zap.String("source", "server"))
	fileManager, err := files.NewFileManager(config.UploaderDir, logger)
	if err != nil {
		return nil, err
	}
	srv := &Server{
		router:      mux.NewRouter(),
		logger:      childLogger,
		config:      config,
		fileManager: *fileManager,
	}
	srv.registerHandlers()
	srv.registerMiddlewares()
	srv.handler = srv.router

	return srv, nil
}
