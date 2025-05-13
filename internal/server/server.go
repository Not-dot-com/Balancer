package server

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// BackendServer представляет бэкенд-сервер
type BackendServer struct {
	URL          *url.URL
	Alive        bool
	mux          sync.RWMutex
	reverseProxy *httputil.ReverseProxy
}

// NewBackendServer создает новый экземпляр BackendServer
func NewBackendServer(url *url.URL) *BackendServer {
	return &BackendServer{
		URL:          url,
		Alive:        true,
		reverseProxy: httputil.NewSingleHostReverseProxy(url),
	}
}

// SetAlive устанавливает состояние сервера (жив/не жив)
func (s *BackendServer) SetAlive(alive bool) {
	s.mux.Lock()
	s.Alive = alive
	s.mux.Unlock()
}

// IsAlive возвращает состояние сервера
func (s *BackendServer) IsAlive() bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	return s.Alive
}

// ServeHTTP перенаправляет запросы на бэкенд-сервер
func (s *BackendServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	s.reverseProxy.ServeHTTP(rw, req)
}
