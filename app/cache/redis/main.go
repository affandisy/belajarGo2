package main

import (
	"log/slog"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/pobyzaarif/go-cache"
	cfg "github.com/pobyzaarif/go-config"
	"github.com/redis/go-redis/v9"
)

type Config struct {
	RedisHost     string `env:"REDIS_HOST"`
	RedisPort     string `env:"REDIS_PORT"`
	RedisPassword string `env:"REDIS_PASSWORD"`
	RedisDB       int    `env:"REDIS_DB"`
}

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

var cacheLoginPrefixKey = "login:"
var listPassword = []string{"correct_Password_123", "wrong_Password_123", "another_wrong_Password_123"}

func main() {
	config := Config{}
	cfg.LoadConfig(config)
	logger.Info("Config loaded")

	loginKey := cacheLoginPrefixKey + "aaa@example.com"

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB:       config.RedisDB,
	})

	redisCache := cache.NewRedisCacheRepository(redisClient)

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
			_ = redisCache.Get(loginKey, &failedLoginCount)

			if failedLoginCount >= 3 {
				logger.Warn("Too many failed login attempt. Please try again later.")
				return
			}

			if listPasswordSelected != listPassword[0] {
				logger.Info("Login failed with password:", slog.String("password", listPasswordSelected))
				count++
				_ = redisCache.Set(loginKey, count, 5*time.Minute)
				return
			}

			logger.Info("Login successfull with password:", slog.String("password", listPasswordSelected))
			count = 0
			redisCache.Delete(loginKey)
		}()
	}
}
