/*
Copyright © 2020 Arnoud Kleinloog

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/akleinloog/lazy-rest/config"

	"github.com/rs/zerolog"
)

var (
	// LogEntryCtxKey is the context.Context key to store the request log entry.
	//LogEntryCtxKey = &contextKey{"LogEntry"}

	// DefaultLogger is called by the Logger middleware handler to log each request.
	// Its made a package-level variable so that it can be reconfigured for custom
	// logging configurations.
	DefaultLogger = New(config.DefaultConfig)
)

// Logger is used for logging.
type Logger struct {
	logger *zerolog.Logger
}

// RequestLogEntry represents an entry in the server's log.
type RequestLogEntry struct {
	Method       string
	URL          string
	UserAgent    string
	Referer      string
	Protocol     string
	RemoteIP     string
	ServerIP     string
	Host         string
	Status       int
	RequestBody  []byte
	ResponseBody []byte
}

// New initializes a new logger
func New(config *config.Config) *Logger {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logLevel := zerolog.InfoLevel

	if config.Debug {
		logLevel = zerolog.DebugLevel
	}

	zerolog.SetGlobalLevel(logLevel)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return &Logger{logger: &logger}
}

// initRequestEntry initializes a new log entry for
func initRequestEntry(request *http.Request) *RequestLogEntry {

	host := request.Host
	if host == "" && request.URL != nil {
		host = request.URL.Host
	}

	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		DefaultLogger.Error().Err(err).Msg("Unable to read request body")
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

	entry := &RequestLogEntry{
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
		entry.ServerIP = ipFromHostPort(localAddress.String())
	}

	return entry
}

// RequestLogger is a middleware that logs the start and end of each request, along
// with some useful data about what was requested, what the response status was,
// and how long it took to return. When standard output is a TTY, Logger will
// print in color, otherwise it will print in black and white. Logger prints a
// request ID if one is provided.
//
// Alternatively, look at https://github.com/goware/httplog for a more in-depth
// http logger with structured logging support.
func RequestLogger(next http.Handler) http.Handler {

	fn := func(writer http.ResponseWriter, request *http.Request) {

		entry := initRequestEntry(request)

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

			DefaultLogger.Info().
				Str("host", entry.Host).
				Str("method", entry.Method).
				Str("url", entry.URL).
				Str("agent", entry.UserAgent).
				Str("referer", entry.Referer).
				Str("protocol", entry.Protocol).
				Str("remoteIp", entry.RemoteIP).
				Str("serverIp", entry.ServerIP).
				Int("status", entry.Status).
				RawJSON("request", entry.RequestBody).
				RawJSON("response", entry.ResponseBody).
				Msg("")
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

// Output duplicates the global logger and sets w as its output.
func (l *Logger) Output(w io.Writer) zerolog.Logger {
	return l.logger.Output(w)
}

// With creates a child logger with the field added to its context.
func (l *Logger) With() zerolog.Context {
	return l.logger.With()
}

// Level creates a child logger with the minimum accepted level set to level.
func (l *Logger) Level(level zerolog.Level) zerolog.Logger {
	return l.logger.Level(level)
}

// Sample returns a logger with the s sampler.
func (l *Logger) Sample(s zerolog.Sampler) zerolog.Logger {
	return l.logger.Sample(s)
}

// Hook returns a logger with the h Hook.
func (l *Logger) Hook(h zerolog.Hook) zerolog.Logger {
	return l.logger.Hook(h)
}

// Debug starts a new message with debug level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Debug() *zerolog.Event {
	return l.logger.Debug()
}

// Info starts a new message with info level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Info() *zerolog.Event {
	return l.logger.Info()
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error starts a new message with error level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Error() *zerolog.Event {
	return l.logger.Error()
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Fatal() *zerolog.Event {
	return l.logger.Fatal()
}

// Panic starts a new message with panic level. The message is also sent
// to the panic function.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Panic() *zerolog.Event {
	return l.logger.Panic()
}

// WithLevel starts a new message with level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) WithLevel(level zerolog.Level) *zerolog.Event {
	return l.logger.WithLevel(level)
}

// Log starts a new message with no level. Setting zerolog.GlobalLevel to
// zerolog.Disabled will still disable events produced by this method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Log() *zerolog.Event {
	return l.logger.Log()
}

// Print sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Print.
func (l *Logger) Print(v ...interface{}) {
	l.logger.Print(v...)
}

// Printf sends a log event using debug level and no extra field.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.logger.Printf(format, v...)
}

// Ctx returns the Logger associated with the ctx. If no logger
// is associated, a disabled logger is returned.
func (l *Logger) Ctx(ctx context.Context) *Logger {
	return &Logger{logger: zerolog.Ctx(ctx)}
}

func IsJSON(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}
