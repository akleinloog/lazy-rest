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
package config

import (
	"github.com/spf13/viper"
)

var (
	// configFile points to the configuration file.
	configFile string

	// DefaultConfig holds the default configuration.
	DefaultConfig = AppConfig()
)

type Config struct {
	Debug  bool `env:"DEBUG, defaults to false"`
	Server serverConf
}

type serverConf struct {
	Port int `env:"LAZYREST_PORT, defaults to 8080"`
}

// AppConfig returns the application configuration.
func AppConfig() *Config {

	var config Config
	config.Debug = viper.GetBool("debug")
	config.Server.Port = viper.GetInt("port")
	if config.Server.Port == 0 {
		config.Server.Port = 8080
	}
	return &config
}
