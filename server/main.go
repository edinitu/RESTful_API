package main

import (
	"flag"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
)

var (
	numberOfInstances int
	port              int
	balancerPort      int
)

const ip string = "127.0.0.1"
const MAX_INSTANCES = 5

func main() {
	flag.IntVar(&numberOfInstances, "no_of_instances", 0, "number of API instances")
	flag.IntVar(&port, "port", -1, "API instance port")
	flag.IntVar(&balancerPort, "balancer_port", 7000, "Load Balancer port")
	flag.Parse()

	if port == -1 {
		log.Fatal("Please input port number for the API instance")
	}
	if numberOfInstances > MAX_INSTANCES {
		log.Fatal("Too many instances, max is 5")
	}

	addr := ip + ":" + strconv.Itoa(port)

	log.Printf("Starting server with address %s\n", addr)

	wco := new(WordsCount)
	wco.Init()

	log.Println("Addresses of other instances are: ", wco.otherAddresses)

	err := rpc.Register(wco)
	if err != nil {
		log.Fatal(err)
	}
	rpc.HandleHTTP()
	err = wco.fileMgr.LoadDataFromFile(wco, port%10)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Current cache contains: ", wco.wordFrequencies)

	healthy := new(Healthy)
	http.Handle("/healthy", healthy)
	http.Handle("/words", wco)
	log.Fatal(http.ListenAndServe(addr, nil))
}
