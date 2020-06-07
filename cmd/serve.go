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
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/akleinloog/lazy-rest/util/logger"
	"github.com/spf13/cobra"
	"io"
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

	key := getURLWithSlashRemovedIfNeeded(request)

	content, prs := memory[key]
	if prs {
		// Request matches a single item, we can return it
		respondWithContent(writer, content)
		return
	}

	var itemsInCollection = make(map[string]interface{})

	for location, item := range memory {
		if strings.HasPrefix(location, key) {
			var subLocation = strings.TrimPrefix(location, key+"/")
			if path.Base(subLocation) == subLocation {
				itemsInCollection[location] = item
			}
		}
	}

	if len(itemsInCollection) > 0 {
		contentItems := make([]interface{}, 0, len(itemsInCollection))
		for _, item := range itemsInCollection {
			contentItems = append(contentItems, item)
		}
		respondWithContent(writer, contentItems)
	} else {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func handlePOST(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashAddedIfNeeded(request)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		server.Logger().Error().Err(err).Msg("Invalid Request Body received")
		http.Error(writer, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// Remove whitespace
	x := bytes.TrimLeft(body, " \t\r\n")

	isArray := len(x) > 0 && x[0] == '['
	// isObject := len(x) > 0 && x[0] == '{'

	decoder := json.NewDecoder(request.Body)

	if isArray {
		// read the array token first
		// then the decoder will iterate over the items in the array
		_, err = decoder.Token()
		if err != nil {
			server.Logger().Error().Err(err).Msg("Error parsing JSON token")
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	var itemsInRequest = make(map[string]interface{})

	for decoder.More() {

		var content map[string]interface{}

		err := decoder.Decode(&content)

		if err != nil {

			server.Logger().Error().Err(err).Msg("Invalid JSON received")
			http.Error(writer, "Invalid JSON", http.StatusBadRequest)
			return

		} else {

			id, prs := content["id"]
			if !prs {
				id = createId()
				content["id"] = id
			}

			contentKey := fmt.Sprintf("%s%s", key, id)
			itemsInRequest[contentKey] = content
		}
	}

	for key, element := range itemsInRequest {
		memory[key] = element
	}

	writer.WriteHeader(http.StatusCreated)
	respond(writer, fmt.Sprintf("Created %d items", len(itemsInRequest)))
}

func handlePUT(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashRemovedIfNeeded(request)

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
			message := fmt.Sprintf("Mismatch between id field `%s` and address `%s`", contentId, resourceId)
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

	key := getURLWithSlashRemovedIfNeeded(request)

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
	if !strings.HasSuffix(key, "/") {
		return key + "/"
	}
	return key
}

func getURLWithSlashRemovedIfNeeded(request *http.Request) string {
	key := request.URL.Path[1:]
	if strings.HasSuffix(key, "/") {
		return strings.TrimSuffix(key, "/")
	}
	return key
}

func createId() string {

	random := make([]byte, 10)
	n, err := io.ReadFull(rand.Reader, random)
	if n != len(random) || err != nil {
		server.Logger().Error().Err(err).Msg("Error while creating if")
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(random)
}
