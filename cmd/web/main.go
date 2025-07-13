package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
    "github.com/patrickmn/go-cache" 

    "github.com/ingwar1991/go_webpage_analyzer/internal/analyzer"
    "github.com/ingwar1991/go_webpage_analyzer/internal/helper"
)

type PageData struct {
    Input string
	Result *analyzer.Result
	Error  string
}


var urlCache *cache.Cache

func main() {
    // default TTL 10min, cleanup every 15min
    urlCache = cache.New(10*time.Minute, 15*time.Minute) 

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Incoming request: %s %s", r.Method, r.URL.Path)

		data := PageData{}

		switch r.Method {
		case http.MethodPost:
            tmpl, err := helper.LoadTemplate("result")
            if err != nil {
                log.Fatalf("Failed to parse templates: %v", err)
            }

			url := helper.NormalizeURL(r.FormValue("url"))
			log.Printf("Analyzing URL: %s", url)

            client := http.Client{
                Timeout: 30 * time.Second,
            }

            resp, err := client.Get(url)
            if err != nil {
                log.Fatalf("Failed to access url: %v", err)
            }
            defer resp.Body.Close()

            if cached, found := urlCache.Get(url); found {
                log.Printf("Cache hit for: %s", url)

                data.Result = cached.(*analyzer.Result)
            } else {
                log.Printf("Cache miss for: %s. Analyzing...", url)

                result, code, err := analyzer.AnalyzePage(url)
                if err != nil {
                    log.Printf("Error analyzing URL %s: %v", url, err)

                    data.Error = err.Error()
                    w.WriteHeader(code)
                } else {
                    data.Result = result
                    urlCache.Set(url, result, cache.DefaultExpiration)
                }
            }

			if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
				log.Printf("Template error: %v", err)
			}

		case http.MethodGet:
            tmpl, err := helper.LoadTemplate("form")
            if err != nil {
                log.Fatalf("Failed to parse templates: %v", err)
            }

			if err := tmpl.ExecuteTemplate(w, "base", data); err != nil {
				log.Printf("Template error: %v", err)

				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Server started on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop // wait for shutdown signal

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited gracefully")
}
