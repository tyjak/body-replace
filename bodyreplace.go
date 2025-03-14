package main

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
)

// Config structure for plugin configuration
type Config struct {}

// CreateConfig initializes the plugin configuration
func CreateConfig() *Config {
	return &Config{}
}

type ReplaceStars struct {
	next   http.Handler
	name   string
	regex  *regexp.Regexp
	replace string
}

// New creates a new ReplaceStars middleware
func New(_ context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	regex := regexp.MustCompile(`\*\*(.*?)\*\*`)
	return &ReplaceStars{
		next:   next,
		name:   name,
		regex:  regex,
		replace: "<b>$1</b>",
	}, nil
}

func (r *ReplaceStars) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	buffer := &bytes.Buffer{}
	responseWriter := &responseWrapper{
		ResponseWriter: rw,
		body:           buffer,
	}

	r.next.ServeHTTP(responseWriter, req)

	modifiedContent := r.regex.ReplaceAllString(buffer.String(), r.replace)
	rw.Header().Set("Content-Length", string(len(modifiedContent)))
	rw.Write([]byte(modifiedContent))
}

type responseWrapper struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (r *responseWrapper) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func (r *responseWrapper) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}

