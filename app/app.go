/*
Copyright © 2020 Arnoud Kleinloog

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
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
