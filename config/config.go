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
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

type Config struct {
}

var cfgFile string

// New returns the configuration.
func New() Config {
	return Config{}
}

func (*Config) Port() int {
	port := viper.GetInt("port")
	if port == 0 {
		port = 8080
	}
	return port
}

func (*Config) InMemory() bool {
	return viper.GetBool("in-memory")
}

// initConfig reads in config file and ENV variables if set.
func Initialize() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".lazy-rest" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lazy-rest")
	}

	viper.SetEnvPrefix("LAZY_REST")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func InitializeRootFlags(rootCmd *cobra.Command) {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.lazy-rest.yaml)")
	rootCmd.PersistentFlags().IntP("port", "p", 0, "port number of the HTTP Server (default is 8080)")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	rootCmd.PersistentFlags().Bool("in-memory", false, "use in memory storage instead of file system")
	viper.BindPFlag("in-memory", rootCmd.PersistentFlags().Lookup("in-memory"))
	rootCmd.PersistentFlags().Lookup("in-memory").NoOptDefVal = "true"
}
