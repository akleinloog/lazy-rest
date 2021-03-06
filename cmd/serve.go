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
package cmd

import (
	"github.com/akleinloog/lazy-rest/pkg/rest"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the REST Server",
	Long: `Starts the REST Server at port 8080.
It will start accepting requests, returning what has been put in.`,
	Run: func(cmd *cobra.Command, args []string) {
		rest.Listen()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
