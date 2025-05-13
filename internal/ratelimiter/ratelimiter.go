package ratelimiter

import (
	"sync"
	"time"
)

// === Rate Limiter (Упрощенная реализация) ===

// RateLimiter struct
type RateLimiter struct {
	requestTimes []time.Time // slice для хранения времени запросов
	mux          sync.Mutex  // Мьютекс для защиты requestTimes
}

// NewRateLimiter конструктор
func NewRateLimiter() RateLimiter {
	return RateLimiter{
		requestTimes: make([]time.Time, 0),
		mux:          sync.Mutex{},
	}
}

// allowRequest проверяет, можно ли выполнить запрос
func (rl *RateLimiter) AllowRequest() bool {
	rl.mux.Lock()
	defer rl.mux.Unlock()

	now := time.Now()

	// Удаляем старые записи
	for i := 0; i < len(rl.requestTimes); i++ {
		if now.Sub(rl.requestTimes[i]) > time.Second {
			rl.requestTimes = rl.requestTimes[i+1:]
			i-- // adjust index after the slice is modified
		}
	}

	// Проверяем лимит
	if len(rl.requestTimes) >= 20 { // Используем BurstSize как ограничение (захардкожено)
		return false
	}

	// Добавляем время нового запроса
	rl.requestTimes = append(rl.requestTimes, now)
	return true
}
