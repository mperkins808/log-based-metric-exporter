package main

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/api/healthy"
	"github.com/mperkins808/log-based-metric-exporter/server/pkg/engine"
	m "github.com/mperkins808/log-based-metric-exporter/server/pkg/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func main() {

	mux := chi.NewRouter()

	// Routes
	mux.Get("/healthy", healthy.Healthy)
	mux.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
		m.ResetPrometheusGauges()
	})

	// ENGINE
	go engine.RuleEngine()

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "10015"
	}

	log.Infof("Starting server on :%v", PORT)
	err := http.ListenAndServe(":"+PORT, mux)
	log.Fatal(err)
}
