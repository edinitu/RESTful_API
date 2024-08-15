package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"strings"
)

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
