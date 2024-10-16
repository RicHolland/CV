package middleware

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"text/template"
	"time"
)

const jsonFormat = `{` +
	`"timestamp":"{{ .Timestamp }}",` +
	`"method":"{{ .Method }}",` +
	`"path":"{{ .Path }}",` +
	`"duration":"{{ .Duration }}",` +
	`"statusCode":{{ .StatusCode }}` +
	"}\n"

type LogData struct {
	Timestamp, Method, Path, Duration string
	StatusCode                        int
}

type wrappedWriter struct {
	writer     http.ResponseWriter
	statusCode int
}

func (w *wrappedWriter) Header() http.Header {
	return w.writer.Header()
}

func (w *wrappedWriter) Write(b []byte) (int, error) {
	return w.writer.Write(b)
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.writer.WriteHeader(statusCode)
	w.statusCode = statusCode
}

type templateLogger struct {
	template *template.Template
	outMu    sync.Mutex
	out      io.Writer
}

func (tl *templateLogger) Execute(data *LogData) error {
	tl.outMu.Lock()
	defer tl.outMu.Unlock()
	return tl.template.Execute(tl.out, data)
}

var jsonTemplate = templateLogger{
	template: template.Must(template.New("JSONLog").Parse(jsonFormat)),
	outMu:    sync.Mutex{},
	out:      os.Stderr,
}

func JSONLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		wrapped := &wrappedWriter{
			writer:     w,
			statusCode: http.StatusOK,
		}

		start := time.Now()

		next.ServeHTTP(wrapped, r)

		end := time.Now()

		l := &LogData{
			Timestamp:  end.Format(time.RFC3339Nano),
			Method:     r.Method,
			Path:       r.URL.Path,
			Duration:   end.Sub(start).String(),
			StatusCode: wrapped.statusCode,
		}
		err := jsonTemplate.Execute(l)
		if err != nil {
			log.Panic(err)
		}
	})
}
