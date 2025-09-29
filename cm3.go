package cm3

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _port = regexp.MustCompile(`:\d+$`)

// Returns X-Real-Ip if it is set, r.RemoteAddr without /:\d+$/ otherwise.
// Removes all occurances of `[` and `]`.
func RemoteIp(r *http.Request) string {
	var ip string
	if realIp, hasRealIp := r.Header["X-Real-Ip"]; hasRealIp && len(realIp) == 1 {
		ip = realIp[0]
	} else {
		ip = _port.ReplaceAllString(r.RemoteAddr, "")
	}
	ip = strings.ReplaceAll(ip, "[", "")
	ip = strings.ReplaceAll(ip, "]", "")
	return ip
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
