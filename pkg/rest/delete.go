package rest

import (
	"github.com/akleinloog/lazy-rest/pkg/storage"
	"net/http"
)

func handleDELETE(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashRemovedIfNeeded(request)

	wasPresent, err := storage.Remove(key)

	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	} else {
		if wasPresent {
			writer.WriteHeader(http.StatusAccepted)
			respond(writer, "")
		} else {
			http.Error(writer, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}
