package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
)

var global_config string

func main() {
	// some virtual init task
	myInitial() // 1. will always be killed for 20s is over initialDelaySeconds=10, but used in /preStart

	go getTooManyMemory()

	//
	// please run with cmd: ./bin/amd64/HTTPServer -logtostderr=true
	//
	flag.Parse()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthzHandler) // 4. livenessProbe
	http.HandleFunc("/preStop", preStopHandler)
	http.HandleFunc("/getData", getDataHandler)
	err := http.ListenAndServe(":8080", nil)
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
	fmt.Println(VERSION)
	glog.Infoln(w.Header())

	// glog.Infoln(r.Method)
	glog.Infoln(ReadUserIP(r))
	glog.Infoln(strings.Split(ReadUserIP(r), ":")[0])

	io.WriteString(w, fmt.Sprintf("hello version is: %s\n", VERSION))
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

// 2. pre stop, /data is a hostpath volume for check
func preStopHandler(w http.ResponseWriter, r *http.Request) {
	write_file("preStop: good night! " + time.Now().String() + "\n")

	w.WriteHeader(200)
	io.WriteString(w, "preStop\n")
	fmt.Println("preStop: good night!", time.Now()) // see log from 'kubectl logs'
}

// 3. it will OOM
func getTooManyMemory() {
	//time.Sleep(time.Duration(8) * time.Second) // it will oom quickly
	time.Sleep(time.Duration(80) * time.Second)
	var mem []int
	for true {
		mem = append(mem, 1)
	}
}

func getDataHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	io.WriteString(w, "our config is: "+global_config+"\n")
}
