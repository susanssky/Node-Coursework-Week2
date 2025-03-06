package api

import (
	"node-week-02-with-chi/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *APIServer) Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.SetHeader("Content-Type", "application/json"))

	messageHandler := handlers.New()

	router.Route("/api/v1/messages", func(r chi.Router) {
		r.Get("/", messageHandler.GetAllMessages)
		r.Post("/", messageHandler.CreateMessage)
		r.Get("/latest", messageHandler.GetLatestMessages)
		r.Get("/search", messageHandler.GetSearchedMessages)

		r.Get("/{messageId}", messageHandler.GetMessage)
		r.Put("/{messageId}", messageHandler.UpdateMessage)
		r.Delete("/{messageId}", messageHandler.DeleteMessage)
	})

	return router
}
