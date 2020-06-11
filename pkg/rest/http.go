package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/akleinloog/lazy-rest/app"
	"github.com/akleinloog/lazy-rest/pkg/logger"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
)

// initRequestLog initializes a new log entry for a request.
func initRequestLog(request *http.Request) *logger.RequestLog {

	host := request.Host
	if host == "" && request.URL != nil {
		host = request.URL.Host
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		app.Log.Error(err, "Unable to read request body")
	} else {

		request.Header.Get("content-type")

		var f interface{}
		err := json.Unmarshal(body, &f)
		if err != nil {
			body = []byte(fmt.Sprintf("%q", body))
		}

		request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	}

	//requestBody := fmt.Sprintf("%q", body)

	logEntry := &logger.RequestLog{
		Host:        host,
		Method:      request.Method,
		URL:         request.URL.String(),
		UserAgent:   request.UserAgent(),
		Referer:     request.Referer(),
		Protocol:    request.Proto,
		RemoteIP:    ipFromHostPort(request.RemoteAddr),
		RequestBody: body,
	}

	if localAddress, ok := request.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		logEntry.ServerIP = ipFromHostPort(localAddress.String())
	}

	return logEntry
}

// RequestLogger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Log will
// print in color, otherwise it will print in black and white. Log prints a
// request ID if one is provided.
//
// Alternatively, look at https://github.com/goware/httplog for a more in-depth
// http-handling logger with structured logging support.
func requestLogger(next http.Handler) http.Handler {

	fn := func(writer http.ResponseWriter, request *http.Request) {

		entry := initRequestLog(request)

		rec := httptest.NewRecorder()

		defer func() {

			entry.Status = rec.Code

			if entry.Status == 0 {
				entry.Status = http.StatusOK
			}

			entry.ResponseBody = rec.Body.Bytes()

			// this copies the recorded response to the response writer
			for k, v := range rec.HeaderMap {
				writer.Header()[k] = v
			}
			writer.WriteHeader(rec.Code)
			rec.Body.WriteTo(writer)

			app.Log.LogRequest(entry)
		}()

		next.ServeHTTP(rec, request)
	}

	return http.HandlerFunc(fn)
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}
	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}
	return h
}

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
	//	server.Log().Error(err, "Error while unmarshalling json")
	//}
	//content := fmt.Sprintf("%s", message)
	// _, err := fmt.Fprint(writer, message)
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(message)
	if err != nil {
		app.Log.Error(err, "Error while responding to request")
	}
}

func respond(writer http.ResponseWriter, message string) {

	_, err := fmt.Fprint(writer, message)
	if err != nil {
		app.Log.Error(err, "Error while responding to request")
	}
}
