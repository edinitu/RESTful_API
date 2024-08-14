package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/rpc"
	"sync"
)

type WordsCount struct {
	lock            sync.Mutex
	wordFrequencies map[string]int

	fileMgr FileMgr
}

func (wco *WordsCount) Init() {
	wco.wordFrequencies = make(map[string]int)
	wco.fileMgr = new(FileMgrImpl)
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
		wco.wordFrequencies[word]++
		update[word] = wco.wordFrequencies[word]
	}
	err := wco.fileMgr.AppendDataToFile(wco, words)
	if err != nil {
		return err
	}

	client, err := rpc.DialHTTP("tcp", "127.0.0.1:7002")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	return client.Call(
		"WordsCount.UpdateCacheForWords",
		&Update{update}, reply)

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
	}
}
