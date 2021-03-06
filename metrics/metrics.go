package metrics

import (
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jordyv/gocam/alerting"
	"github.com/jordyv/gocam/config"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	interval                = 5 * time.Second
	currentAlertImagesGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gocam_current_alert_images",
		Help: "Number of alert images currently in storage",
	})
	alertGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "gocam_alerts",
		Help: "Number of alerts",
	})
)

func init() {
	prometheus.MustRegister(currentAlertImagesGauge)
	prometheus.MustRegister(alertGauge)
}

// Server is the metrics server struct
type Server struct {
	config *config.Config
}

// New initiates a new metrics server
func New(config *config.Config) *Server {
	return &Server{config: config}
}

func (s *Server) refreshMetrics() {
	files, _ := ioutil.ReadDir(s.config.AlertImagePath)
	currentAlertImagesGauge.Set(float64(len(files)))

}

// Notify gets called by the alert manager and will update the metrics
func (s Server) Notify(alert *alerting.Alert) {
	alertGauge.Add(1.00)
}

// Listen starts listening at the metrics endpoint
func (s *Server) Listen() {
	go func() {
		for {
			s.refreshMetrics()
			time.Sleep(interval)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(s.config.MetricsAddr, nil))
}
