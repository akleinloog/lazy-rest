package app

import (
	"github.com/akleinloog/lazy-rest/config"
	"github.com/akleinloog/lazy-rest/util/logger"
)

// Server represents the overall application.
type Server struct {
	logger *logger.Logger
	config *config.Config
}

// Instance creates a new Server with config and logger..
func Instance() *Server {
	return &Server{logger: logger.New(config.AppConfig()), config: config.AppConfig()}
}

// Logger provides access to the global logger.
func (app *Server) Logger() *logger.Logger {
	return app.logger
}

// Config provides access to the global configuration.
func (app *Server) Config() *config.Config {
	return app.config
}
