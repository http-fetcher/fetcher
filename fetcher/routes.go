package fetcher

func (s *server) routes() {
	s.router.Post("/api/fetcher", s.handleUpdateCreate())
	s.router.Delete("/api/fetcher/{id}", s.handleDelete())
	s.router.Get("/api/fetcher", s.handleList())
	s.router.Get("/api/fetcher/{id}/history", s.handleHistory())
}
