package rest

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/akleinloog/lazy-rest/pkg/storage"
	"io"
	"io/ioutil"
	"net/http"
)

func handlePOST(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashAddedIfNeeded(request)

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		app.Log.Error(err, "Invalid Request Body received")
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
			app.Log.Error(err, "Error parsing JSON token")
			http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	var itemsInRequest = make(map[string]interface{})

	for decoder.More() {

		var content map[string]interface{}

		err := decoder.Decode(&content)

		if err != nil {

			app.Log.Error(err, "Invalid JSON received")
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
		storage.Store(key, element)
	}

	writer.WriteHeader(http.StatusCreated)
	respond(writer, fmt.Sprintf("Created %d items", len(itemsInRequest)))
}

func createId() string {

	random := make([]byte, 10)
	n, err := io.ReadFull(rand.Reader, random)
	if n != len(random) || err != nil {
		app.Log.Error(err, "Error while creating if")
		panic(err)
	}

	return base64.RawURLEncoding.EncodeToString(random)
}
