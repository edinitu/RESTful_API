package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"strconv"
	"sync"
)

type Balancer interface {
	AddApiInstance(api API)
	GetApiInstances() []API
	Rotate() API
	GetNextValidPeer() API
}

type BalancerImpl struct {
	APIInstances []API      // list with all API instances in our system
	lock         sync.Mutex // locks the rotation of instances and the current variable
	current      int        // holds the index to the current API instance we want to direct the request to
}

// Reverse proxy implementation that directs the request to another server
func (b *BalancerImpl) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	api := b.GetNextValidPeer()
	if api == nil {
		res.WriteHeader(http.StatusInternalServerError)
		_, err := res.Write([]byte("Service unavailable"))
		if err != nil {
			log.Printf("ERROR: %s\n", err.Error())
		}
		return
	}

	url := api.GetUrl()
	log.Printf("Directing request to the following URL: %s\n", url.String())
	proxy := httputil.NewSingleHostReverseProxy(url)
	proxy.ServeHTTP(res, req)
}

func (b *BalancerImpl) AddApiInstance(api API) {
	b.APIInstances = append(b.APIInstances, api)
}

func (b *BalancerImpl) GetApiInstances() []API {
	return b.APIInstances
}

func (b *BalancerImpl) Rotate() API {
	b.lock.Lock()
	b.current = (b.current + 1) % len(b.GetApiInstances())
	b.lock.Unlock()
	return b.APIInstances[b.current]
}

func (b *BalancerImpl) GetNextValidPeer() API {
	for i := 0; i < len(b.GetApiInstances()); i++ {
		nextPeer := b.Rotate()
		if nextPeer.IsAlive() {
			return nextPeer
		} else {
			log.Println("WARN: Instance with URL is not available: ", nextPeer.GetUrl().String())
		}
	}
	return nil
}

func populateBalancer(balancer Balancer) {
	log.Println("Adding instances to load balancer")
	for i := 0; i < numberOfInstances; i++ {
		api := new(APIImpl)
		apiSocket := ip + ":" + strconv.Itoa(port+i+1)
		err := api.SetUrl(apiSocket)
		if err != nil {
			log.Println("ERROR: Could not set url: ", apiSocket)
			continue
		}
		err = api.SetHealthUrl(apiSocket)
		if err != nil {
			log.Println("ERROR: Could not set health url: ", apiSocket)
			continue
		}
		balancer.AddApiInstance(api)
		log.Println("Added: ", api.URL.String())
	}
	if len(balancer.GetApiInstances()) == 0 {
		log.Fatal("Could not add any api instance")
	}
}
