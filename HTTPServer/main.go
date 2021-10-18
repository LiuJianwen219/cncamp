package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/golang/glog"
)

func main() {
	//
	// please run with cmd: ./bin/amd64/HTTPServer -logtostderr=true
	//
	flag.Parse()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthzHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	// 1. wirte request Headers to response's Header
	for k, v := range r.Header {
		for _, sv := range v {
			w.Header().Add(k, sv)
		}
	}

	// 2. set environment parameter to response's Header
	VERSION := os.Getenv("VERSION")
	w.Header().Add("VERSION", VERSION)
	fmt.Println(VERSION)
	glog.Infoln(w.Header())

	// 3. get IP from request
	// glog.Infoln(r.Method)
	glog.Infoln(ReadUserIP(r))
	glog.Infoln(strings.Split(ReadUserIP(r), ":")[0])

	io.WriteString(w, fmt.Sprintf("hello version is: %s\n", VERSION))
}

// 4. return 200 for path /healthz
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "healthz\n")
}

func ReadUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}
