package rest

import (
	"github.com/akleinloog/lazy-rest/pkg/storage"
	"net/http"
)

func handleGET(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashRemovedIfNeeded(request)

	content, exists, err := storage.Retrieve(key)
	if err != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if exists {
		// Request matches a single item, we can return it
		respondWithContent(writer, content)
		return
	}

	var itemsInCollection, getErr = storage.RetrieveCollection(key)

	if getErr != nil {
		http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
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
