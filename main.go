package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"time"
	"log"
	"context"
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"encoding/json"
)

type Server struct {
	ID     string    `json:"id"`
	Name   string `json:"name"`
	RemoteID string `json:"remote_id"`
	Type string `json:"type"`
	Status string `json:"status"`
	IP string `json:"ip"`
}

type RemoteServer struct {
	ID     string    `json:"id"`
	Name   string `json:"name"`
	Type string `json:"type"`
	Status string `json:"status"`
	IP string `json:"ip"`
}

type Sandbox struct {
	ID string `json:"id"`
	Status string `json:"status"`
	PreviewURL string `json:"preview_url"`
	WebsocketURL string `json:"websocket_url"`
}

type SandboxCreateRequest struct {
	ImageName string `json:"image_name"`
	StartCommand string `json:"start_command,omitempty"`
}

type ContainerCreateRequest struct {
	ImageName string `json:"image_name"`
	StartCommand string `json:"start_command,omitempty"`
}

type ServiceCreateRequest struct {
	Name string `json:"name"`
}

type Container struct {
	ID string `json:"id"`
	ServiceID string `json:"service_id"`
	ServerID string `json:"server_id"`
	SandboxID string `json:"sandbox_id"`
	Host string `json:"host"`
	ImageName string `json:"image_name"`
	StartCommand string `json:"start_command",omitempty`
	Status string `json:"status"`
}

var serviceHandler = NewServiceHandler()
var serverHandler = NewServerHandler()
var serverPort = "8002"

func main() {
	err := godotenv.Load()
    if err != nil {
        log.Println("Warning: Error loading .env file")
    }

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/api", func(r chi.Router) {
		r.Route("/services", func(r chi.Router) {
			r.Get("/", func(w http.ResponseWriter, r *http.Request) {
				result, err := serviceHandler.ListServices()
				if err != nil {
					returnErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				data := &ServiceCreateRequest{}
	
				if err := json.NewDecoder(r.Body).Decode(data); err != nil {
					returnErrorResponse(w, "Invalid request body", http.StatusBadRequest)
					return
				}
				result, err := serviceHandler.CreateService(data.Name)
				if err != nil {
					returnErrorResponse(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Get("/{serviceID}", func(w http.ResponseWriter, r *http.Request) {
				serviceID := chi.URLParam(r, "serviceID")
				result, err := serviceHandler.GetService(serviceID)
				if err != nil {
					returnErrorResponse(w, "Not found", http.StatusNotFound)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(result)
			})

			r.Delete("/{serviceID}", func(w http.ResponseWriter, r *http.Request) {
				serviceID := chi.URLParam(r, "serviceID")
				err := serviceHandler.DeleteService(serviceID)
				if err != nil {
					returnErrorResponse(w, "Not found", http.StatusNotFound)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			})

			r.Route("/{serviceID}/containers", func(r chi.Router) {
				r.Get("/", func(w http.ResponseWriter, r *http.Request) {
					serviceID := chi.URLParam(r, "serviceID")
					service, err := serviceHandler.GetService(serviceID)
					if err != nil {
						returnErrorResponse(w, "Not found", http.StatusNotFound)
						return
					}
					result, err := service.ListContainers()
					if err != nil {
						returnErrorResponse(w, "Internal server error", http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(result)
				})

				r.Post("/", func(w http.ResponseWriter, r *http.Request) {
					serviceID := chi.URLParam(r, "serviceID")
					service, err := serviceHandler.GetService(serviceID)
					if err != nil {
						returnErrorResponse(w, "Not found", http.StatusNotFound)
						return
					}

					data := &ContainerCreateRequest{}
	
					if err := json.NewDecoder(r.Body).Decode(data); err != nil {
						returnErrorResponse(w, "Invalid request body", http.StatusBadRequest)
						return
					}

					container, err := service.CreateContainer(data.ImageName, data.StartCommand)
					if err != nil {
						returnErrorResponse(w, "Internal Server error", http.StatusInternalServerError)
						return
					}

					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(container)
				})
			})
		})
	})

	server := &http.Server{
		Addr: fmt.Sprintf(":%s", serverPort),
		Handler: r,
	}

	// Channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Starting server on port", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Shutdown signal received, shutting down gracefully...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server shutdown complete")
}

func returnErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]string{
		"message": message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
