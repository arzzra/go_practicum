package api

import (
	"errors"
	"github.com/arzzra/go_practicum/internal/server"
	"github.com/arzzra/go_practicum/internal/server/metric"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

type Handler struct {
	Server *server.Server
	Router *chi.Mux
}

func MakeHandler(srv *server.Server) (*Handler, error) {
	if srv == nil {
		return nil, errors.New("server is not initialized")
	}

	router := chi.NewRouter()
	if router == nil {
		return nil, errors.New("router is not initialized")
	}

	h := &Handler{
		Server: srv,
		Router: router,
	}

	h.initMiddleware()
	h.initRouting()

	return h, nil
}

func (h *Handler) initMiddleware() {
	h.Router.Use(middleware.Recoverer)
}

func withMetricTypeValidator(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		metricType := chi.URLParam(r, "Type")
		if err := metric.MetricType(metricType).Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusNotImplemented)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(f)
}

func (h *Handler) initRouting() {
	h.Router.Route("/update/{Type}/{Name}/{Value}",
		func(r chi.Router) {
			r.Use(withMetricTypeValidator)
			r.Post("/", h.postMetric)
		})

	h.Router.Route("/value/{Type}/{Name}",
		func(r chi.Router) {
			r.Use(withMetricTypeValidator)
			r.Get("/", h.getMetric)
		})

	h.Router.Get("/", h.getAllMetrics)
}
