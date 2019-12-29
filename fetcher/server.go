package fetcher

import (
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	maxBodySize int64
	router *chi.Mux
}

func NewServer(maxBodySize int64) *server {
	s := &server{maxBodySize: maxBodySize}
	s.router = chi.NewRouter()
	s.routes()
	return s
}

// make server an http.Handler
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBodySize)
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
			if err.Error() == "http: request body too large" {
				http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		log.Printf("Input: %v", in)
		out := output{Id: 1}
		s.respond(w, r, out, http.StatusOK)
	}
}

type IdHandlerFunc func(http.ResponseWriter, *http.Request, int)

func withId(h IdHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h(w, r, id)
	}
}

func (s *server) handleDelete() http.HandlerFunc {
	return withId(func(w http.ResponseWriter, r *http.Request, id int) {
		log.Printf("Deleting id: %d", id)
	})
}

func (s *server) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}

func (s *server) handleHistory() http.HandlerFunc {
	return withId(func(w http.ResponseWriter, r *http.Request, id int) {
		log.Printf("History for id: %d", id)
	})
}
