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
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

// Logger is used for logging.
type Logger struct {
	logger *zerolog.Logger
}

// RequestLog is used to log http requests and their responses.
type RequestLog struct {
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
func New() Logger {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logLevel := zerolog.InfoLevel

	zerolog.SetGlobalLevel(logLevel)

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	return Logger{logger: &logger}
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

// LogRequest writes the info contained in a RequestLog as Info to the log stream.
func (l *Logger) LogRequest(log *RequestLog) {
	l.logger.Info().
		Str("host", log.Host).
		Str("method", log.Method).
		Str("url", log.URL).
		Str("agent", log.UserAgent).
		Str("referer", log.Referer).
		Str("protocol", log.Protocol).
		Str("remoteIp", log.RemoteIP).
		Str("serverIp", log.ServerIP).
		Int("status", log.Status).
		RawJSON("request", log.RequestBody).
		RawJSON("response", log.ResponseBody).
		Msg("")
}

// Warn starts a new message with warn level.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Warn() *zerolog.Event {
	return l.logger.Warn()
}

// Error Logs a new message with error level.
func (l *Logger) Error(err error, msg string) {
	l.logger.Error().Err(err).Msg(msg)
}

// Fatal starts a new message with fatal level. The os.Exit(1) function
// is called by the Msg method.
//
// You must call Msg on the returned event in order to send the event.
func (l *Logger) Fatal(err error, msg string) {
	l.logger.Fatal().Err(err).Msg(msg)
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
