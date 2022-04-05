package api

import (
	"fmt"
	"github.com/arzzra/go_practicum/internal/server/metric"
	"github.com/go-chi/chi/v5"
	"html/template"
	"net/http"
)

const templateHTML = `
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Metrics}}<div>{{ .Name }}: {{ .StringValue }}</div>{{end}}
	</body>
	</html>`

func (h *Handler) postMetric(wrt http.ResponseWriter, req *http.Request) {
	typeM := metric.MetricType(chi.URLParam(req, "Type"))
	name := chi.URLParam(req, "Name")
	value := chi.URLParam(req, "Value")

	m, err := metric.MakeMetricStruct(name, typeM, value)
	if err != nil {
		http.Error(wrt, err.Error(), http.StatusBadRequest)
		return
	}
	h.Server.Storage.SaveMetric(m)
	wrt.WriteHeader(http.StatusOK)
}

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	typeM := metric.MetricType(chi.URLParam(r, "Type"))
	name := chi.URLParam(r, "Name")
	m, err := h.Server.Storage.GetMetricFromStorage(typeM, name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if m == nil {
		http.Error(w, fmt.Sprintf("Metric %s not found", name), http.StatusNotFound)
		return
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, m.GetValueString())
}

type metricsForHTML struct {
	Title   string
	Metrics []metric.Metric
}

func (h *Handler) getAllMetrics(w http.ResponseWriter, r *http.Request) {

	t, err := template.New("getAllMetric").Parse(templateHTML)
	if err != nil {
		errHTTP := http.StatusInternalServerError
		http.Error(w, err.Error(), errHTTP)
		return
	}
	typeContent := r.Header.Get("content-type")
	m, err := h.Server.Storage.GetAllMetricFromStorage()
	if err != nil {
		errCode := http.StatusInternalServerError
		http.Error(w, err.Error(), errCode)
		return
	}

	data := metricsForHTML{
		Title:   "Metrics List",
		Metrics: *m,
	}
	w.Header().Set("content-type", typeContent)
	w.WriteHeader(http.StatusOK)
	_ = t.Execute(w, data)
}
