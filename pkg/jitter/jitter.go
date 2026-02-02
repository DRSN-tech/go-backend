// Package jitter предоставляет утилиты для добавления случайности в интервалы отступления (backoff),
// чтобы предотвратить эффект «буйного стада» (thundering herd) в распределённых системах.
package jitter

import (
	"math/rand"
	"sync"
	"time"
)

// DefaultJitter — стандартный коэффициент джиттера (50%)
const DefaultJitter = 0.5

var (
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	randMutex  sync.Mutex
)

// Duration возвращает продолжительность с применённым джиттером.
// Результат находится в диапазоне [d, d*(1+jitterFactor)].
func Duration(d time.Duration, jitterFactor float64) time.Duration {
	randMutex.Lock()
	jitter := globalRand.Float64() * jitterFactor * float64(d)
	randMutex.Unlock()
	return d + time.Duration(jitter)
}

// DurationWithSeed возвращает продолжительность с джиттером, используя заданный генератор случайных чисел.
// Полезно для тестирования или когда требуется детерминированное поведение.
func DurationWithSeed(d time.Duration, jitterFactor float64, rng *rand.Rand) time.Duration {
	return d + time.Duration(rng.Float64()*jitterFactor*float64(d))
}

// ExponentialBackoff вычисляет экспоненциальное отступление с джиттером.
// base — начальная длительность отступления,
// max — максимальная длительность отступления,
// attempt — номер текущей попытки повтора (нумерация с нуля),
// jitterFactor — коэффициент джиттера (например, 0.5 означает ±50%).
func ExponentialBackoff(base, max time.Duration, attempt int, jitterFactor float64) time.Duration {
	backoff := base
	for i := 0; i < attempt; i++ {
		backoff *= 2
		if backoff > max {
			backoff = max
			break
		}
	}
	return Duration(backoff, jitterFactor)
}
