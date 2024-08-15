package main

import (
	"log"
	"net/http"
	"net/url"
	"sync"
)

type API interface {
	IsAlive() bool
	SetUrl(socket string) error
	SetHealthUrl(socket string) error
	GetUrl() *url.URL
}

type APIImpl struct {
	URL       *url.URL
	healthURL *url.URL
	lock      sync.Mutex
}

func (api *APIImpl) IsAlive() bool {
	api.lock.Lock()
	defer api.lock.Unlock()
	resp, err := http.Get(api.healthURL.String())
	if err != nil {
		log.Println("ERROR: Could not check alive status of api instance")
		return false
	}
	if resp.Status == "200 OK" {
		return true
	}
	return false
}

func (api *APIImpl) SetUrl(socket string) error {
	apiUrl, err := url.Parse("http://" + socket)
	if err != nil {
		return err
	}
	api.URL = apiUrl
	return nil
}

func (api *APIImpl) SetHealthUrl(socket string) error {
	apiUrl, err := url.Parse("http://" + socket + "/healthy")
	if err != nil {
		return err
	}
	api.healthURL = apiUrl
	return nil
}

func (api *APIImpl) GetUrl() *url.URL {
	return api.URL
}
