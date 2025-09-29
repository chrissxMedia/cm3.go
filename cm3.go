package cm3

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _port = regexp.MustCompile(`:\d+$`)

// Returns X-Real-Ip if it is set, r.RemoteAddr without /:\d+$/ otherwise.
func RemoteIp(r *http.Request) string {
	if realIp, hasRealIp := r.Header["X-Real-Ip"]; hasRealIp && len(realIp) == 1 {
		return realIp[0]
	} else {
		return _port.ReplaceAllString(r.RemoteAddr, "")
	}
}

// Registers the given prometheus metrics and serves them at /metrics.
func HandleMetrics(metrics ...prometheus.Collector) {
	for _, m := range metrics {
		prometheus.MustRegister(m)
	}
	// TODO: think about logging
	http.Handle("/metrics", promhttp.Handler())
}

// Registers an http handler using http.HandleFunc and logs all requests.
func HandleFunc(location string, handler func(w http.ResponseWriter, r *http.Request)) {
	if handler != nil {
		http.HandleFunc(location, func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s %s from %s (%s) to %s\n",
				r.Method, r.URL.Path, r.Proto, RemoteIp(r), r.UserAgent(), r.Host)
			handler(w, r)
		})
	}
}

func ListenAndServeHttp(addr string, handler func(w http.ResponseWriter, r *http.Request)) {
	HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe(addr, nil))
}

func ListenAndServeHttps(addr string, certFile string, keyFile string, handler func(w http.ResponseWriter, r *http.Request)) {
	HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServeTLS(addr, certFile, keyFile, nil))
}
