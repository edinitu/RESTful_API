package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"strconv"
	"strings"
)

var (
	instanceId   int
	balancerPort int
)

const ip string = "127.0.0.1"

type TextPostBody struct {
	Text string
}

type WordsGetBody struct {
	Words []string
}

type Update struct {
	Words map[string]int
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

func main() {
	flag.IntVar(&instanceId, "instance_id", 0, "server instance id")
	flag.IntVar(&balancerPort, "balancer_port", 7000, "load balancer port")
	flag.Parse()

	serverPort := balancerPort + 1 + instanceId%5
	addr := ip + ":" + strconv.Itoa(serverPort)

	log.Printf("Starting server with address %s\n", addr)

	wco := new(WordsCount)
	wco.Init()
	err := rpc.Register(wco)
	if err != nil {
		log.Fatal(err)
	}
	rpc.HandleHTTP()
	err = wco.fileMgr.LoadDataFromFile(wco, instanceId)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(wco.wordFrequencies)
	http.Handle("/words", wco)
	log.Fatal(http.ListenAndServe(addr, nil))
}
