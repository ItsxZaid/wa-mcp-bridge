package server

import (
	"fmt"
	"net/http"
	"wa-mcp-bridge/internal/store"
	"wa-mcp-bridge/internal/whatsapp"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	router *chi.Mux
	addr   string
	store  store.Store
	waBot  *whatsapp.Bot
}

func New(httpPort string, store store.Store, waBot *whatsapp.Bot) *http.Server {
	s := &Server{
		router: chi.NewRouter(),
		addr:   fmt.Sprintf(":%s", httpPort),
		store:  store,
		waBot:  waBot,
	}

	s.router.Use(middleware.Logger)

	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/whatsapp", s.waBot.RegisterRoutes)
	})

	return &http.Server{
		Addr:    s.addr,
		Handler: s.router,
	}
}
