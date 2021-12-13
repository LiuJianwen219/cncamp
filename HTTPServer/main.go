package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LiuJianwen/cncamp/HTTPServer/metrics"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var global_config string

var sub_version = "v1.4"

func main() {
	myInitial()

	go getTooManyMemory()

	//
	// please run with cmd: ./bin/amd64/HTTPServer -logtostderr=true
	//
	flag.Parse()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthzHandler)
	http.HandleFunc("/preStop", preStopHandler)
	http.HandleFunc("/getData", getDataHandler)

	metrics.Register()
	http.HandleFunc("/count", countHandler)
	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		for _, sv := range v {
			w.Header().Add(k, sv)
		}
	}

	VERSION := os.Getenv("VERSION")
	w.Header().Add("VERSION", VERSION)
	w.Header().Add("sub_version", sub_version)
	fmt.Println(VERSION)
	glog.Infoln(w.Header())

	// glog.Infoln(r.Method)
	glog.Infoln(ReadUserIP(r))
	glog.Infoln(strings.Split(ReadUserIP(r), ":")[0])

	io.WriteString(w, fmt.Sprintf("hello version is: %s, sub_version is: %s\n", VERSION, sub_version))
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "healthz\n")
	user := r.Header.Get("user")
	fmt.Println("healthz", user, time.Now()) // see log from 'kubectl logs'
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

func myInitial() {
	fmt.Println("init start", time.Now())
	time.Sleep(time.Duration(5) * time.Second)
	// time.Sleep(time.Duration(20) * time.Second) // 2. will always be killed for 20s is over initialDelaySeconds=10
	fmt.Println("init over", time.Now())

	// 5. get config from volume, file
	bytes, err := ioutil.ReadFile("/data/config.ini")
	if err != nil {
		fmt.Printf("Unable to read file: %v\n", err)
	}
	global_config = string(bytes)
	write_file("config: we get config " + global_config)
	fmt.Println(global_config)

	err = ioutil.WriteFile("/tmp/healthy", []byte("I can read now."), 0755)
	if err != nil {
		fmt.Printf("Unable to write file: %v\n", err)
	}
}

func preStopHandler(w http.ResponseWriter, r *http.Request) {
	write_file("preStop: good night! " + time.Now().String() + "\n")

	w.WriteHeader(200)
	io.WriteString(w, "preStop\n")
	fmt.Println("preStop: good night!", time.Now()) // see log from 'kubectl logs'
}

func getTooManyMemory() {
	//time.Sleep(time.Duration(8) * time.Second) // it will oom quickly
	time.Sleep(time.Duration(8000) * time.Second)
	var mem []int
	for true {
		mem = append(mem, 1)
	}
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "our config is: "+global_config+"\n")
}

func countHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("count metrics")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(0, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "======== Details of the http request header: ========\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	glog.V(4).Infof("Respond in %d ms", delay)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}
