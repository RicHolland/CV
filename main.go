package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"cv/components"
	"cv/middleware"

	"github.com/a-h/templ"
)

const (
	Version = "v0.1.0"
	banner  = `
 ££££££\  ££\    ££\
££  __££\ ££ |   ££ |
££ /  \__|££ |   ££ |
££ |      \££\  ££  |
££ |       \££\££  /
££ |  ££\   \£££  /
\££££££  |   \£  /
 \______/     \_/ version %s, built with %s

`
)

//go:embed all:assets
var assets embed.FS

// Render 404 Error page
func error404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	components.Error404(r.URL.Path).Render(context.Background(), w)
}

func main() {
	host := flag.String("host", "localhost", "Host")
	port := flag.Int("port", 8080, "Port to serve on")
	flag.Parse()

	router := http.NewServeMux()

	// Serve static content from `./assets/`
	router.Handle("GET /assets/", http.FileServerFS(assets))

	// Endpoints
	router.Handle("GET /{$}", templ.Handler(components.Index()))

	// Handle 404s
	router.HandleFunc("/", error404)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", *host, *port),
		Handler: middleware.AttachMiddleware(router, middleware.JSONLogging),
	}

	fmt.Printf(banner, Version, runtime.Version())
	log.Printf("Listening and serving on: %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
