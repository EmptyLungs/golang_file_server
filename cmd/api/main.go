package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EmptyLungs/golang_file_server/pkg/api"
	"github.com/EmptyLungs/golang_file_server/pkg/files"
	"github.com/EmptyLungs/golang_file_server/pkg/grpc"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	fs := pflag.NewFlagSet("default", pflag.ContinueOnError)
	fs.String("host", "127.0.0.1", "Host to bind HTTP server to")
	fs.Int("port", 8080, "Port to bind HTTP server to")
	fs.String("upload-dir", "./data", "Files directory")
	fs.Int64("upload-max-file-size", 50, "Upload file size limit")
	fs.String("level", "info", "Log level")
	fs.Duration("http-server-timeout", 30*time.Second, "server read and write timeout duration")
	fs.Int("grpc-port", 0, "gRPC port")
	fs.String("grpc-service-name", "gofs", "gRPC service name")

	err := fs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		fmt.Fprintf(os.Stderr, "Error: %s\n\n", err.Error())
		fs.PrintDefaults()
		os.Exit(2)
	}
	viper.BindPFlags(fs)
	viper.SetEnvPrefix("FS_SERVER")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	var srvCfg api.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		logger.Panic("HTTP Server config unmarshal error", zap.Error(err))
	}
	dirFs := os.DirFS(srvCfg.UploaderDir)
	fileManager, err := files.NewFileManager(srvCfg.UploaderDir, dirFs, logger)
	if err != nil {
		logger.Panic(err.Error())
	}

	var grpcCfg grpc.Config
	if err := viper.Unmarshal(&grpcCfg); err != nil {
		logger.Panic("gRPC Server config unmarshal error", zap.Error(err))
	}

	if grpcCfg.Port > 0 {
		grpcSrv, _ := grpc.NewServer(&grpcCfg, logger)
		grpcSrv.ListenAndServe()
	}
	srv, err := api.NewServer(&srvCfg, logger, *fileManager)
	if err != nil {
		logger.Panic("server_error", zap.Error(err))
	}
	srv.ListenAndServe()
}
