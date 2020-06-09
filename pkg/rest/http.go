package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

func respond(writer http.ResponseWriter, message string) {

	_, err := fmt.Fprint(writer, message)
	if err != nil {
		server.Logger().Error().Err(err).Msg("Error while responding to request")
	}
}
