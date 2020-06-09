package rest

import (
	"github.com/akleinloog/lazy-rest/pkg/storage"
	"net/http"
)

func handleDELETE(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashRemovedIfNeeded(request)

	wasPresent := storage.Remove(key)
	if wasPresent {
		writer.WriteHeader(http.StatusAccepted)
		respond(writer, "")
	} else {
		http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}
