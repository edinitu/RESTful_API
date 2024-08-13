package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var (
	instanceId   int
	balancerPort int

	FILE_PATH string
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

func persistDataToFile(AllData *WordsCount, words []string) error {
	var textToWrite string

	for _, word := range words {
		textToWrite += word + ":" + strconv.Itoa(AllData.wordFrequencies[word])
		textToWrite += "\n"
	}
	f, err := os.OpenFile(FILE_PATH, os.O_WRONLY|os.O_APPEND, 0660)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("ERROR: Could not close file")
		}
	}(f)
	if err != nil {
		return err
	}
	_, err = f.WriteString(textToWrite)
	return err
}

func (wco *WordsCount) UpdateCacheAndPersist(words []string) error {
	for _, word := range words {
		wco.wordFrequencies[word]++
	}
	return persistDataToFile(wco, words)
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

		err = wco.UpdateCacheAndPersist(words)
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

func main() {
	flag.IntVar(&instanceId, "instance_id", 0, "server instance id")
	flag.IntVar(&balancerPort, "balancer_port", 7000, "load balancer port")
	flag.Parse()

	serverPort := balancerPort + 1 + instanceId%5
	addr := ip + ":" + strconv.Itoa(serverPort)

	log.Printf("Starting server with address %s\n", addr)

	wco := new(WordsCount)
	wco.Init()

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	FILE_PATH = dir + "/words_instance_" + strconv.Itoa(instanceId) + ".txt"
	log.Println("Open file " + FILE_PATH)
	f, err := os.OpenFile(FILE_PATH, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	data, err := os.ReadFile(FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reload cache with file data")
	strData := strings.TrimRight(string(data), "\n")
	rows := strings.Split(strData, "\n")
	for _, row := range rows {
		pair := strings.Split(row, ":")
		if len(pair) != 2 {
			log.Printf("Warning, row has more than 2 elements: %v, skipping", row)
			continue
		}
		i, err := strconv.Atoi(pair[1])
		if err != nil {
			log.Printf("Warning, could not parse frequency for row: %v, skipping", row)
			continue
		}
		wco.wordFrequencies[pair[0]] = i
	}

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(wco.wordFrequencies)
	http.Handle("/words", wco)
	log.Fatal(http.ListenAndServe(addr, nil))
}
