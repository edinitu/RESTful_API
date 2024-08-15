package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var FILE_PATH string

type FileMgr interface {
	LoadDataFromFile(AllData *WordsCount, instanceId int) error
	AppendDataToFile(AllData *WordsCount, words []string) error
}

type FileMgrImpl struct{}

func (fmi *FileMgrImpl) LoadDataFromFile(AllData *WordsCount, instanceId int) error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return err
	}

	FILE_PATH = dir + "/words_instance_" + strconv.Itoa(instanceId) + ".txt"
	log.Println("Open file " + FILE_PATH)
	f, err := os.OpenFile(FILE_PATH, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Reload cache with file data")
	strData := strings.Trim(string(data), "\n ")
	rows := strings.Split(strData, "\n")
	for _, row := range rows {
		if row == "" {
			continue
		}
		pair := strings.Split(row, ":")
		if len(pair) != 2 {
			log.Printf("Warning, row should have precisely 2 elements: %v, skipping", row)
			continue
		}
		i, err := strconv.Atoi(pair[1])
		if err != nil {
			log.Printf("Warning, could not parse frequency for row: %v, skipping", row)
			continue
		}
		AllData.wordFrequencies[pair[0]] = i
	}

	return f.Close()
}

func (fmi *FileMgrImpl) AppendDataToFile(AllData *WordsCount, words []string) error {
	var textToWrite string

	for _, word := range words {
		normalizedWord := toLowerAndStripSpecialChars(word)
		textToWrite += normalizedWord + ":" + strconv.Itoa(AllData.wordFrequencies[normalizedWord])
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
