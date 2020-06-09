package rest

import (
	"encoding/json"
	"fmt"
	"github.com/akleinloog/lazy-rest/pkg/storage"
	"io/ioutil"
	"net/http"
	"path"
)

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
	storage.Store(key, content)

	writer.WriteHeader(http.StatusAccepted)
	respond(writer, "")
}