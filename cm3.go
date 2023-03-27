package cm3

import (
	"fmt"
	"log"
	"net/http"
)

func HandleFunc(location string, handler func(w http.ResponseWriter, r *http.Request)) {
	if handler != nil {
		http.HandleFunc(location, func(w http.ResponseWriter, r *http.Request) {
			log.Println(fmt.Sprintf("%s %s %s from %s (%s) to %s",
				r.Method, r.URL.Path, r.Proto, r.RemoteAddr, r.UserAgent(), r.Host))
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
