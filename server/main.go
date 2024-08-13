package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
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

type WordsGetBody struct {
	Words []string
}

func (wco *WordsCount) Init() {
	wco.wordFrequencies = make(map[string]int)
}

func decodeAndValidatePost(text *TextPostBody, body io.ReadCloser) ([]string, error) {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&text)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return nil, err
	}

	words := strings.Fields(text.Text)
	if len(words) == 0 {
		return nil, errors.New("could not update memory. " +
			"Cause is either empty text or unknown json format.\n accepted json is: \n" +
			"{\"Text\" : \"this is an example\"'}")
	}

	return words, nil
}

func decodeAndValidateGet(words *WordsGetBody, body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&words)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return err
	}

	if len(words.Words) == 0 {
		return errors.New("please provide at least one word")
	}

	return nil
}

func (wco *WordsCount) UpdateCacheAndPersist(words []string) {
	for _, word := range words {
		wco.wordFrequencies[word]++
	}
	//TODO persist to file
}

func (wco *WordsCount) GetResponseFromCache(words WordsGetBody) map[string]int {
	response := make(map[string]int)

	for _, word := range words.Words {
		if freq, ok := wco.wordFrequencies[word]; ok {
			response[word] = freq
		} else {
			response[word] = 0
		}
	}

	return response
}

func (wco *WordsCount) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	wco.lock.Lock()
	defer wco.lock.Unlock()
	if req.Method == "POST" {
		var text TextPostBody

		words, err := decodeAndValidatePost(&text, req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Printf("POST: ERROR %v", err.Error())
			_, err = rw.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
			return
		}

		wco.UpdateCacheAndPersist(words)

		rw.WriteHeader(http.StatusOK)
		log.Printf("\nPOST:\nUpdated map for text: %s\n", text.Text)
		return
	}

	if req.Method == "GET" {
		var words WordsGetBody

		err := decodeAndValidateGet(&words, req.Body)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			log.Printf("GET: ERROR %v", err.Error())
			_, err = rw.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
			return
		}

		response := wco.GetResponseFromCache(words)
		
		rw.WriteHeader(http.StatusOK)
		err = json.NewEncoder(rw).Encode(response)
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
		}
		log.Printf("\nGET:\n words %v.\n Response:\n %v", words.Words, response)
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
