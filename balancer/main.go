package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
)

var (
	numberOfInstances int
	port              int
)

const ip string = "127.0.0.1"

type Balancer struct {
	addresses []string
}

func (b *Balancer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Println("Directing request to the following URL: http://127.0.0.1:7001/")
	serverUrl, _ := url.Parse("http://127.0.0.1:7001/")
	proxy := httputil.NewSingleHostReverseProxy(serverUrl)
	proxy.ServeHTTP(res, req)
}

func main() {
	flag.IntVar(&numberOfInstances, "number_of_instances", 0, "number of API instances")
	flag.IntVar(&port, "port", 7000, "load balancer port")
	flag.Parse()

	addr := ip + ":" + strconv.Itoa(port)

	log.Printf("Starting balancer server with address %s\n", addr)

	balancer := new(Balancer)
	http.Handle("/", balancer)
	log.Fatal(http.ListenAndServe(addr, nil))
}
