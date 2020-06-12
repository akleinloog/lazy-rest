package rest

import (
	"encoding/json"
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/akleinloog/lazy-rest/pkg/storage"
	"io/ioutil"
	"net/http"
	"path"
)

func handlePUT(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashRemovedIfNeeded(request)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		app.Log.Error(err, "Unable to read request body")
		writer.WriteHeader(http.StatusBadRequest)
		respond(writer, "")
		return
	}

	// TODO: Check Content Type / Support multiple content types
	// memory[key] = body
	// contentType := request.Header.Retrieve("content-type")

	var content interface{}
	err = json.Unmarshal(body, &content)
	if err != nil {
		// Invalid JSON
		app.Log.Error(err, "Invalid JSON received")
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
	err = storage.Store(key, content)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusAccepted)
	respond(writer, "")
}
