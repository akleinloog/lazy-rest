package rest

import (
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/akleinloog/lazy-rest/pkg/util/logger"
	"net/http"
	"os"
)

var (
	host   string = "unknown"
	server *app.Server
)

func Listen() {

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
