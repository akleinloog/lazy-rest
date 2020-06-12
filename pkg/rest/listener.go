package rest

import (
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"net/http"
	"os"
)

var (
	host = "unknown"
)

func Listen() {

	currentHost, err := os.Hostname()
	if err != nil {
		app.Log.Info().Msgf("Could not determine host name:", err)
	} else {
		host = currentHost
	}

	app.Log.Info().Msgf("Starting Lazy REST Server on " + host)

	requestHandler := http.HandlerFunc(HandleRequest)

	http.Handle("/", requestLogger(requestHandler))

	address := fmt.Sprintf("%s:%d", "", app.Config.Port())

	err = http.ListenAndServe(address, nil)
	if err != nil {
		app.Log.Fatal(err, "Error while listening for requests")
	}
}

// HandleRequest determines the appropriate action to take based on the http method.
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
