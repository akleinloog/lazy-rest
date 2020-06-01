/*
Copyright © 2020 Arnoud Kleinloog <arnoud@kleinloog.ch>

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
package main

import (
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/akleinloog/lazy-rest/util/logger"
	"net/http"
	"os"
)

var (
	requestNr int64  = 0
	host      string = "unknown"
)

func main() {

	server := app.Instance()

	currentHost, err := os.Hostname()
	if err != nil {
		server.Logger().Info().Msgf("Could not determine host name:", err)
	} else {
		host = currentHost
	}

	server.Logger().Info().Msgf("Starting Lazy REST Server on " + host)

	requestHandler := http.HandlerFunc(HandleRequest)

	http.Handle("/", logger.RequestLogger(requestHandler))

	address := fmt.Sprintf("%s:%d", "", server.Config().Server.Port)

	err = http.ListenAndServe(address, nil)
	if err != nil {
		server.Logger().Fatal().Err(err)
	}
}

// Hello gives out a simple hello message
func HandleRequest(w http.ResponseWriter, r *http.Request) {

	requestNr++
	message := fmt.Sprintf("Go Hello %d from %s on %s ./%s\n", requestNr, host, r.Method, r.URL.Path[1:])
	fmt.Fprint(w, message)
}
