package main

import (
	"io"
	"strings"
	"testing"
)

func TestPostDecoding(t *testing.T) {
	{
		stringReader := strings.NewReader("{\"Text\": \"This is a sample text\"}")
		stringReadCloser := io.NopCloser(stringReader)

		expectedWords := []string{"This", "is", "a", "sample", "text"}
		var text *TextPostBody
		actualWords, err := decodeAndValidatePost(text, stringReadCloser)
		if err != nil {
			t.Fatalf("Got error while decoding")
		}
		for idx, _ := range expectedWords {
			if expectedWords[idx] != actualWords[idx] {
				t.Fatalf("Expected %s, got %s at index %d", expectedWords[idx], actualWords[idx], idx)
			}
		}
	}
	{
		stringReader := strings.NewReader("{\"Text\": \"This is   a sample   text with   spaces\"}")
		stringReadCloser := io.NopCloser(stringReader)

		expectedWords := []string{"This", "is", "a", "sample", "text", "with", "spaces"}
		var text TextPostBody
		actualWords, err := decodeAndValidatePost(&text, stringReadCloser)
		if err != nil {
			t.Fatalf("Got error while decoding")
		}
		for idx, _ := range expectedWords {
			if expectedWords[idx] != actualWords[idx] {
				t.Fatalf("Expected %s, got %s at index %d", expectedWords[idx], actualWords[idx], idx)
			}
		}
	}
}

func TestGetDecoding(t *testing.T) {
	{
		stringReader := strings.NewReader("{\"Words\": [\"test_word\"]}")
		stringReadCloser := io.NopCloser(stringReader)

		expectedWords := []string{"test_word"}
		var actualWords WordsGetBody
		err := decodeAndValidateGet(&actualWords, stringReadCloser)
		if err != nil {
			t.Fatalf("Got error while decoding")
		}
		if len(actualWords.Words) == 0 {
			t.Fatalf("Got empty words list")
		}
		for idx, _ := range expectedWords {
			if expectedWords[idx] != actualWords.Words[idx] {
				t.Fatalf("Expected %s, got %s at index %d",
					expectedWords[idx], actualWords.Words[idx], idx)
			}
		}
	}
	{
		stringReader := strings.NewReader("{\"Words\": [\"test_word\", \"test_word2\", \"test_word3\"]}")
		stringReadCloser := io.NopCloser(stringReader)

		expectedWords := []string{"test_word", "test_word2", "test_word3"}
		var actualWords WordsGetBody
		err := decodeAndValidateGet(&actualWords, stringReadCloser)
		if err != nil {
			t.Fatalf("Got error while decoding")
		}
		if len(actualWords.Words) == 0 {
			t.Fatalf("Got empty words list")
		}
		for idx, _ := range expectedWords {
			if expectedWords[idx] != actualWords.Words[idx] {
				t.Fatalf("Expected %s, got %s at index %d",
					expectedWords[idx], actualWords.Words[idx], idx)
			}
		}
	}
}
