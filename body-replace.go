package bodyreplace

import (
	"bytes"
	"io"
	"net/http"
	"regexp"
)

type Config struct {
	Search  string `json:"search,omitempty"`
	Replace string `json:"replace,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type BodyReplace struct {
	next    http.Handler
	name    string
	search  *regexp.Regexp
	replace string
}

func New(next http.Handler, config *Config, name string) (http.Handler, error) {
	re, err := regexp.Compile(config.Search)
	if err != nil {
		return nil, err
	}

	return &BodyReplace{
		next:    next,
		name:    name,
		search:  re,
		replace: config.Replace,
	}, nil
}

func (br *BodyReplace) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	buffer := &bytes.Buffer{}
	writer := &responseWriter{ResponseWriter: rw, buffer: buffer}

	br.next.ServeHTTP(writer, req)

	modifiedBody := br.search.ReplaceAllString(buffer.String(), br.replace)
	rw.WriteHeader(writer.statusCode)
	_, _ = rw.Write([]byte(modifiedBody))
}

type responseWriter struct {
	http.ResponseWriter
	buffer     *bytes.Buffer
	statusCode int
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	return rw.buffer.Write(b)
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}

