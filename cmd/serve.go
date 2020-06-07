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
	"encoding/json"
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/akleinloog/lazy-rest/util/logger"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	host   string                 = "unknown"
	memory map[string]interface{} = make(map[string]interface{})
	server *app.Server
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the REST Server",
	Long: `Starts the REST Server at port 8080.
It will start accepting requests, returning what has been put in.`,
	Run: func(cmd *cobra.Command, args []string) {
		listen()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func listen() {

	server = app.Instance()

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
func HandleRequest(writer http.ResponseWriter, request *http.Request) {

	switch request.Method {
	case "GET":
		handleGET(writer, request)
	case "POST":
		handlePOST(writer, request)
	case "PUT":
		handlePUT(writer, request)
	case "DELETE":
		handleDELETE(writer, request)
	default:
		http.Error(writer, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	}
}

func handleGET(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashAddedIfNeeded(request)

	content, prs := memory[key]
	if prs {
		respondWithContent(writer, content)
	} else {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func handlePOST(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashRemovedIfNeeded(request)

	decoder := json.NewDecoder(request.Body)

	var content map[string]interface{}

	err := decoder.Decode(&content)
	if err != nil {
		// Invalid JSON
		http.Error(writer, "Invalid JSON", http.StatusBadRequest)
	} else {

		id, prs := content["id"]

		if prs {
			contentKey := fmt.Sprintf("%s%s", key, id)
			memory[contentKey] = content
			writer.WriteHeader(http.StatusCreated)
			respond(writer, "")
		} else {
			http.Error(writer, "Missing id field", http.StatusBadRequest)
		}
	}
}

func handlePUT(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashAddedIfNeeded(request)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		server.Logger().Error().Err(err).Msg("Unable to read request body")
		writer.WriteHeader(http.StatusBadRequest)
		respond(writer, "")
		return
	}

	// TODO: Check Content Type / Support multiple content types
	// memory[key] = body
	// contentType := request.Header.Get("content-type")

	var content interface{}
	err = json.Unmarshal(body, &content)
	if err != nil {
		// Invalid JSON
		server.Logger().Error().Err(err).Msg("Invalid JSON received")
		http.Error(writer, "Invalid JSON", http.StatusBadRequest)
		return
	}

	jsonContent := content.(map[string]interface{})

	resourceId := path.Base(key)

	contentId, prs := jsonContent["id"]
	if prs {
		if contentId != resourceId {
			message := fmt.Sprintf("Mismatch between id %s and address %s", contentId, resourceId)
			http.Error(writer, message, http.StatusBadRequest)
			return
		}
	} else {
		jsonContent["id"] = resourceId
	}

	// Valid JSON
	memory[key] = content
	writer.WriteHeader(http.StatusAccepted)
	respond(writer, "")
}

func handleDELETE(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashAddedIfNeeded(request)

	_, prs := memory[key]
	if prs {
		delete(memory, key)
		writer.WriteHeader(http.StatusAccepted)
		respond(writer, "")
	} else {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func respond(writer http.ResponseWriter, message string) {

	_, err := fmt.Fprint(writer, message)
	if err != nil {
		server.Logger().Error().Err(err).Msg("Error while responding to request")
	}
}

func respondWithContent(writer http.ResponseWriter, message interface{}) {

	//content, err := json.Marshal(message)
	//if err != nil {
	//	server.Logger().Error().Err(err).Msg("Error while unmarshalling json")
	//}
	//content := fmt.Sprintf("%s", message)
	// _, err := fmt.Fprint(writer, message)
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(message)
	if err != nil {
		server.Logger().Error().Err(err).Msg("Error while responding to request")
	}
}

func getURLWithSlashAddedIfNeeded(request *http.Request) string {
	key := request.URL.Path[1:]
	if strings.HasSuffix(key, "/") {
		return strings.TrimSuffix(key, "/")
	}
	return key
}

func getURLWithSlashRemovedIfNeeded(request *http.Request) string {
	key := request.URL.Path[1:]
	if !strings.HasSuffix(key, "/") {
		return key + "/"
	}
	return key
}
