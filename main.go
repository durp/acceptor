package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	s := store{
		m: make(map[string]any),
	}
	r.Get("/health", ping)
	r.Route("/api/v1/accept", func(r chi.Router) {
		r.Post("/", s.accept)
		r.Put("/", s.accept)
		r.Get("/{id}", s.get)
	})

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "80"
	}

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + httpPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		fmt.Println(err.Error())
	}
}

func ping(w http.ResponseWriter, _ *http.Request) {
	_ = prettyPrint(w, pong{Status: "ok"})
}

type pong struct {
	Status string `json:"status"`
}

type store struct {
	m map[string]any
}

func (s store) accept(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&s.m)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = prettyPrint(w, &s.m)
}

func (s store) get(w http.ResponseWriter, _ *http.Request) {
	_ = prettyPrint(w, &s.m)
}

func prettyPrint(w io.Writer, v interface{}) (err error) {
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.SetEscapeHTML(false)
	return e.Encode(v)
}
