package fetcher

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"log"
	"net/http"
	"strconv"
)

type server struct {
	maxBodySize int64
	router *chi.Mux
	crawler Crawler
}

func NewServer(maxBodySize int64, router *chi.Mux, crawler Crawler) *server {
	s := &server{maxBodySize: maxBodySize}
	s.crawler = crawler
	s.router = router
	s.routes()
	return s
}

func (s *server) routes() {
	s.router.Post("/api/fetcher", s.handlePut())
	s.router.Delete("/api/fetcher/{Id}", s.handleDelete())
	s.router.Get("/api/fetcher", s.handleList())
	s.router.Get("/api/fetcher/{Id}/history", s.handleHistory())
}

// Make server an http.Handler.
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, s.maxBodySize)
	s.router.ServeHTTP(w, r)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

type IdHandlerFunc func(http.ResponseWriter, *http.Request, int64)

func withId(h IdHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "Id"), 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		h(w, r, id)
	}
}

func (s *server) handlePut() http.HandlerFunc {
	type output struct {
		Id int64 `json:"id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var in Spec
		err := s.decode(w, r, &in)
		if err != nil {
			// TODO: figure out better way to check error type
			if err.Error() == "http: request body too large" {
				http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
			} else {
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
			return
		}
		log.Printf("Put input: %v", in)
		spec := s.crawler.Put(in)
		s.respond(w, r, output{Id: spec.Id}, http.StatusOK)
	}
}

func (s *server) handleDelete() http.HandlerFunc {
	return withId(func(w http.ResponseWriter, r *http.Request, id int64) {
		log.Printf("Delete id: %d", id)
		err := s.crawler.Del(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
	})
}

func (s *server) handleList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("List")
		specs := s.crawler.GetSpecs()
		s.respond(w, r, specs, http.StatusOK)
	}
}

func (s *server) handleHistory() http.HandlerFunc {
	return withId(func(w http.ResponseWriter, r *http.Request, id int64) {
		log.Printf("History id: %d", id)
		results, err := s.crawler.GetResults(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		s.respond(w, r, results, http.StatusOK)
	})
}
