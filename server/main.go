package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var (
	instanceId   int
	balancerPort int
)

const ip string = "127.0.0.1"

type WordsCount struct {
	lock            sync.Mutex
	wordFrequencies map[string]int
}

type TextPostBody struct {
	Text string
}

func (wco *WordsCount) Init() {
	wco.wordFrequencies = make(map[string]int)
}

func (wco *WordsCount) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wco.lock.Lock()
	defer wco.lock.Unlock()
	if req.Method == "POST" {
		decoder := json.NewDecoder(req.Body)
		var text TextPostBody
		err := decoder.Decode(&text)
		if err != nil {
			log.Printf("ERROR: %v\n", err)
			rw.WriteHeader(http.StatusBadRequest)
			_, err := rw.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
			return
		}

		words := strings.Fields(text.Text)
		if len(words) == 0 {
			rw.WriteHeader(http.StatusBadRequest)
			_, err := rw.Write([]byte("empty text. Please provide at least one word"))
			if err != nil {
				panic(err)
			}
			return
		}

		for _, word := range words {
			wco.wordFrequencies[word]++
		}

		rw.WriteHeader(http.StatusOK)
		log.Printf("\nPOST:\nUpdated map for text: %s\n", text.Text)
	}
}

func main() {
	flag.IntVar(&instanceId, "instance_id", 0, "server instance id")
	flag.IntVar(&balancerPort, "balancer_port", 7000, "load balancer port")
	flag.Parse()

	serverPort := balancerPort + 1 + instanceId%5
	addr := ip + ":" + strconv.Itoa(serverPort)

	log.Printf("Starting server with address %s\n", addr)

	wco := new(WordsCount)
	wco.Init()
	http.Handle("/words", wco)
	log.Fatal(http.ListenAndServe(addr, nil))
}
