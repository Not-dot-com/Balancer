package balancer

import (
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"loadbalancer/internal/server"
)

// ServerPool управляет набором бэкенд-серверов
type ServerPool struct {
	servers []*server.BackendServer
	current int
	mux     sync.Mutex
}

// AddServer добавляет новый бэкенд-сервер в пул
func (s *ServerPool) AddServer(server *server.BackendServer) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.servers = append(s.servers, server)
}

// NextServer возвращает следующий доступный сервер, используя round-robin
func (s *ServerPool) NextServer() *server.BackendServer {
	s.mux.Lock()
	defer s.mux.Unlock()

	next := s.current
	for i := 0; i < len(s.servers); i++ {
		server := s.servers[next]
		next = (next + 1) % len(s.servers)
		if server.IsAlive() {
			s.current = next
			return server
		}
	}

	return nil
}

// HealthCheck периодически проверяет состояние бэкенд-серверов
func (s *ServerPool) HealthCheck() {
	for {
		time.Sleep(5 * time.Second) // Упрощенный health check

		for _, server := range s.servers {
			status := "up"
			alive := isBackendAlive(server.URL)
			server.SetAlive(alive)
			if !alive {
				status = "down"
			}
			log.Printf("%s [%s]\n", server.URL, status)
		}
	}
}

// isBackendAlive проверяет, жив ли бэкенд-сервер
func isBackendAlive(url *url.URL) bool {
	resp, err := http.Get(url.String())
	if err != nil {
		log.Println("Ошибка при проверке:", err) // Упрощенная обработка ошибок
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}
