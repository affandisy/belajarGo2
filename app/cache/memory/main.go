package main

import (
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/pobyzaarif/go-cache"
)

var cacheLoginPrefixKey = "login:"
var listPassword = []string{"correct_Password_123", "wrong_Password_123", "another_wrong_Password_123"}

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

func main() {
	loginKey := cacheLoginPrefixKey + "aaa@example.com"

	memoryCache, _ := cache.NewMemoryARCCacheRepository(1)

	var mu sync.Mutex
	count := 0
	for i := 0; i < 10; i++ {
		randSource := rand.New(rand.NewSource(time.Now().UnixNano()))
		indexPassword := randSource.Intn(3)
		listPasswordSelected := listPassword[indexPassword]

		func() {
			mu.Lock()
			defer mu.Unlock()

			time.Sleep(200 * time.Millisecond)

			var failedLoginCount int
			_ = memoryCache.Get(loginKey, &failedLoginCount)

			if failedLoginCount >= 3 {
				logger.Warn("Too many failed login attempts. Please try again later.", slog.String("loginKey", loginKey))
				return
			}

			if listPasswordSelected != listPassword[0] {
				logger.Info("Login failed with password:", slog.String("loginKey", loginKey), slog.String("password", listPasswordSelected))
				count++
				_ = memoryCache.Set(loginKey, count, 5*time.Minute)
				return
			}

			logger.Info("Login successfull with password:", slog.String("loginKey", loginKey), slog.String("password", listPasswordSelected))
			count = 0
			memoryCache.Delete(loginKey)
		}()
	}
}
