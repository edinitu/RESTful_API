package main

//
//type Backend interface {
//	IsAlive() bool
//	GetURL() *url.URL
//	ServeHTTP(http.ResponseWriter, *http.Request)
//}
//
//type API struct {
//	url          *url.URL
//	alive        bool
//	lock         sync.RWMutex
//	reverseProxy *httputil.ReverseProxy
//}
//
//func (api *API) GetURL() *url.URL {
//	return api.url
//}
//
//type roundRobin struct {
//	apis    []API
//	lock    sync.RWMutex
//	current int
//}
//
//func (s *roundRobin) Rotate() API {
//	s.lock.Lock()
//	s.current = (s.current + 1) % s.GetServerPoolSize()
//	s.lock.Unlock()
//	return s.apis[s.current]
//}
//
//func (s *roundRobin) GetNextValidPeer() API {
//	for i := 0; i < s.GetServerPoolSize(); i++ {
//		nextPeer := s.Rotate()
//		if nextPeer.IsAlive() {
//			return nextPeer
//		}
//	}
//	return nil
//}
