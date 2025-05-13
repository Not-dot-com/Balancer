package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"loadbalancer/internal/balancer"
	"loadbalancer/internal/ratelimiter"
	"loadbalancer/internal/server"
)

// === Настраиваемые параметры ===
const (
	ServerPort      = 8080 // Порт для балансировщика нагрузки
	NumberOfServers = 3    // Количество бэкенд-серверов
)

var serverPool balancer.ServerPool
var rateLimiter ratelimiter.RateLimiter // Используем RateLimiter из пакета

func loadBalancer(w http.ResponseWriter, r *http.Request) {
	if !rateLimiter.AllowRequest() {
		w.WriteHeader(http.StatusTooManyRequests)
		fmt.Fprintln(w, "Too Many Requests")
		return
	}

	server := serverPool.NextServer()
	if server != nil {
		server.ServeHTTP(w, r)
		return
	}
	http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
}

func main() {
	rateLimiter = ratelimiter.NewRateLimiter() // Инициализируем rate limiter

	backends := make([]*url.URL, NumberOfServers)
	for i := 0; i < NumberOfServers; i++ {
		backendURL, err := url.Parse(fmt.Sprintf("http://localhost:%d", 8081+i))
		if err != nil {
			log.Fatal(err) // Упрощенная обработка ошибок
		}
		backends[i] = backendURL
	}

	for _, backendURL := range backends {
		server := server.NewBackendServer(backendURL)
		serverPool.AddServer(server)
		log.Println("Configured backend server:", backendURL)
	}

	go serverPool.HealthCheck()

	http.HandleFunc("/", loadBalancer)
	fmt.Println("Load Balancer started on port", ServerPort)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", ServerPort), nil)) // Упрощенная обработка ошибок
}
