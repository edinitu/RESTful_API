package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

var (
	numberOfInstances int
	port              int
)

const ip string = "127.0.0.1"

func main() {
	flag.IntVar(&numberOfInstances, "number_of_instances", 0, "number of API instances")
	flag.IntVar(&port, "port", 7000, "load balancer port")
	flag.Parse()

	if numberOfInstances > 5 {
		log.Fatal("Too many instances. Please input a number lower or equal with 5")
	}

	if numberOfInstances <= 0 {
		log.Fatal("Negative number of instances")
	}

	addr := ip + ":" + strconv.Itoa(port)

	log.Printf("Starting balancer server with address %s\n", addr)

	balancer := new(BalancerImpl)
	populateBalancer(balancer)

	http.Handle("/", balancer)
	log.Fatal(http.ListenAndServe(addr, nil))
}
