package fetcher

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	router *chi.Mux
}

func NewServer() *server {
	s := &server{}
	s.router = chi.NewRouter()
	s.routes()
	return s
}

// make server an http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) handleUpdateCreate() http.HandlerFunc {
	type input struct {
		Url string
		Interval int
	}
	type output struct {
		Id int
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var in input
		err := s.decode(w, r, &in)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Input: %v", in)
		out := output{Id: 1}
		s.respond(w, r, out, http.StatusOK)
	}
}

func (s *server) handleDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("Deleting id: %d", id)
	}
}

func (s *server) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (s *server) handleHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("History for id: %d", id)
	}
}
