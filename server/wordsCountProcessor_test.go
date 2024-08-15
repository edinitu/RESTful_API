package main

import "testing"

func TestInit(t *testing.T) {
	numberOfInstances = 2
	balancerPort = 7000
	port = 7001
	wco := new(WordsCount)
	wco.Init()
	expectedOtherAddresses := []string{"127.0.0.1:7002"}
	if wco.otherAddresses == nil || len(wco.otherAddresses) != 1 ||
		wco.otherAddresses[0] != expectedOtherAddresses[0] {
		t.Fatalf("Expected %v, got %v", expectedOtherAddresses, wco.otherAddresses)
	}
}

func TestToLowerAndTrimWord(t *testing.T) {
	{
		s1 := "Test"
		actual := toLowerAndStripSpecialChars(s1)
		expected := "test"
		if actual != expected {
			t.Fatalf("Expected %s, actual %s", expected, actual)
		}
	}
	{
		s1 := "Test???"
		actual := toLowerAndStripSpecialChars(s1)
		expected := "test"
		if actual != expected {
			t.Fatalf("Expected %s, actual %s", expected, actual)
		}
	}
	{
		s1 := "^&*Test???"
		actual := toLowerAndStripSpecialChars(s1)
		expected := "test"
		if actual != expected {
			t.Fatalf("Expected %s, actual %s", expected, actual)
		}
	}
	{
		s1 := "^&*Test"
		actual := toLowerAndStripSpecialChars(s1)
		expected := "test"
		if actual != expected {
			t.Fatalf("Expected %s, actual %s", expected, actual)
		}
	}
}

func TestGetResponseFromCache(t *testing.T) {
	{
		wco := new(WordsCount)
		wco.Init()
		wco.wordFrequencies["test"] = 3
		wco.wordFrequencies["test2"] = 2

		actual := wco.GetResponseFromCache(WordsGetBody{[]string{"test"}})
		if _, ok := actual["test"]; !ok {
			t.Fatalf("Expected 'test' to be in resulted map, but is not")
		}

		if actual["test"] != 3 {
			t.Fatalf("Expected frequency 3, got %d", actual["test"])
		}

		actual = wco.GetResponseFromCache(WordsGetBody{[]string{"test2"}})
		if _, ok := actual["test2"]; !ok {
			t.Fatalf("Expected 'test2' to be in resulted map, but is not")
		}

		if actual["test2"] != 2 {
			t.Fatalf("Expected frequency 2, got %d", actual["test2"])
		}
	}
	{
		wco := new(WordsCount)
		wco.Init()
		wco.wordFrequencies["test"] = 3
		wco.wordFrequencies["test2"] = 2
		wco.wordFrequencies["test3"] = 1

		actual := wco.GetResponseFromCache(WordsGetBody{[]string{"test", "test2", "test3"}})
		if _, ok := actual["test"]; !ok {
			t.Fatalf("Expected 'test' to be in resulted map, but is not")
		}
		if _, ok := actual["test2"]; !ok {
			t.Fatalf("Expected 'test2' to be in resulted map, but is not")
		}
		if _, ok := actual["test3"]; !ok {
			t.Fatalf("Expected 'test3' to be in resulted map, but is not")
		}

		if actual["test"] != 3 {
			t.Fatalf("Expected frequency 3, got %d", actual["test"])
		}
		if actual["test2"] != 2 {
			t.Fatalf("Expected frequency 2, got %d", actual["test2"])
		}
		if actual["test3"] != 1 {
			t.Fatalf("Expected frequency 1, got %d", actual["test3"])
		}
	}
}
