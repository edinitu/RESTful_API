package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
	"sync"
)

type Healthy struct{}

func (h *Healthy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusOK)
}

type WordsCount struct {
	lock            sync.Mutex
	wordFrequencies map[string]int
	rpcClients      map[string]*rpc.Client
	otherAddresses  []string
	fileMgr         FileMgr
}

func (wco *WordsCount) Init() {
	wco.wordFrequencies = make(map[string]int)
	wco.rpcClients = make(map[string]*rpc.Client)
	wco.fileMgr = new(FileMgrImpl)

	for i := 0; i < numberOfInstances; i++ {
		if i+1 == port%10 {
			continue
		}
		wco.otherAddresses = append(wco.otherAddresses, ip+":"+strconv.Itoa(balancerPort+i+1))
	}
}

func (wco *WordsCount) UpdateCacheForWords(update Update, reply *int) error {
	log.Println("Update received from other API instance")
	var wordsToPersist []string
	for key, val := range update.Words {
		wco.wordFrequencies[key] = val
		wordsToPersist = append(wordsToPersist, key)
	}
	return wco.fileMgr.AppendDataToFile(wco, wordsToPersist)
}

func (wco *WordsCount) UpdateCacheAndPersist(words []string, reply *int) error {
	update := make(map[string]int)
	for _, word := range words {
		normalizedWord := toLowerAndStripSpecialChars(word)
		wco.wordFrequencies[normalizedWord]++
		update[normalizedWord] = wco.wordFrequencies[normalizedWord]
	}
	err := wco.fileMgr.AppendDataToFile(wco, words)
	if err != nil {
		return err
	}

	for _, address := range wco.otherAddresses {
		var client *rpc.Client
		if _, ok := wco.rpcClients[address]; !ok {
			client, err = rpc.DialHTTP("tcp", address)
			if err != nil {
				log.Println("Error dialing: ", err)
				continue
			}
			wco.rpcClients[address] = client
		} else {
			client = wco.rpcClients[address]
		}

		err = client.Call(
			"WordsCount.UpdateCacheForWords",
			&Update{update}, reply)
		if err != nil {
			log.Println("Error while synchronizing other instances, ", err.Error())
		}
	}

	return nil
}

func (wco *WordsCount) GetResponseFromCache(words WordsGetBody) map[string]int {
	response := make(map[string]int)

	for _, word := range words.Words {
		if freq, ok := wco.wordFrequencies[toLowerAndStripSpecialChars(word)]; ok {
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

		var reply int = 0
		err = wco.UpdateCacheAndPersist(words, &reply)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			log.Printf("POST: ERROR %v", err.Error())
			_, err = rw.Write([]byte(err.Error()))
		}

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
		return
	}

	rw.WriteHeader(http.StatusBadRequest)
	_, err := rw.Write([]byte("Method not supported yet"))
	if err != nil {
		panic(err)
	}
}

func toLowerAndStripSpecialChars(s string) string {
	// consider just the string with no special character or digits before or after it
	s = strings.Trim(s, "!@#$%^&*()_{}:<>?=[];',.//|\\0123456789-\n ")
	return strings.ToLower(s)
}
