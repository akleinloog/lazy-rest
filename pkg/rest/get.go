package rest

import (
	"github.com/akleinloog/lazy-rest/pkg/storage"
	"net/http"
)

func handleGET(writer http.ResponseWriter, request *http.Request) {

	key := getURLWithSlashRemovedIfNeeded(request)

	content, prs := storage.Get(key)
	if prs {
		// Request matches a single item, we can return it
		respondWithContent(writer, content)
		return
	}

	var itemsInCollection = storage.GetCollection(key)

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
