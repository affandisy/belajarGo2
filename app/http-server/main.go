package main

import (
	invCtrl "belajarGo2/app/http-server/controller/inventory"
	invRepo "belajarGo2/repository/inventory"
	invSvc "belajarGo2/service/inventory"
	"belajarGo2/util/database"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/julienschmidt/httprouter"
	cfg "github.com/pobyzaarif/go-config"
)

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

type Config struct {
	// AppPort string `env:"APP_PORT"`
	AppVersion string `env:"APP_VERSION"`
	AppHost    string `env:"APP_HOST"`
	AppPort    string `env:"APP_PORT_HTTP_SERVER"`

	DBDriver        string `env:"DB_DRIVER"`
	DBMySQLHost     string `env:"DB_MYSQL_HOST"`
	DBMySQLPort     string `env:"DB_MYSQL_PORT"`
	DBMySQLUser     string `env:"DB_MYSQL_USER"`
	DBMySQLPassword string `env:"DB_MYSQL_PASSWORD"`
	DBMySQLName     string `env:"DB_MYSQL_NAME"`
}

func main() {
	spew.Dump()

	config := Config{}
	cfg.LoadConfig(&config)
	logger.Info("Config loaded")

	// Init db connection
	databaseConfig := database.Config{
		DBDriver:        config.DBDriver,
		DBMySQLHost:     config.DBMySQLHost,
		DBMySQLPort:     config.DBMySQLPort,
		DBMySQLUser:     config.DBMySQLUser,
		DBMySQLPassword: config.DBMySQLPassword,
		DBMySQLName:     config.DBMySQLName,
	}

	// _ = databaseConfig.GetDatabaseConnection()
	db := databaseConfig.GetDatabaseConnection()
	logger.Info("Database client connected!")

	// Dependency Injection
	inventoryRepo := invRepo.NewGormRepository(db)
	inventorySvc := invSvc.NewService(inventoryRepo)
	inventoryCtrl := invCtrl.NewController(logger, inventorySvc)

	// Setup router
	router := httprouter.New()

	router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": config.AppVersion})
	})

	router.GET("/ping", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	})

	// router.GET("/querypath/:id", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	id := p.ByName("id")
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	_ = json.NewEncoder(w).Encode(map[string]string{"message": id})
	// })

	// router.GET("/queryparam", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// 	queryParam := r.URL.Query()

	// 	a := map[string][]string{}
	// 	for k, v := range queryParam {
	// 		a[k] = v
	// 	}
	// 	spew.Dump(a["id"])
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)
	// 	_ = json.NewEncoder(w).Encode(map[string]string{"message": "pong"})
	// })

	// router.POST("/urlencoded", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)

	// 	_ = json.NewEncoder(w).Encode(map[string]string{"message": "b"})
	// })

	// router.POST("/body", func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// 	decoder := json.NewDecoder(r.Body)
	// 	type user struct {
	// 		Name string `json:"name" validate:"required"`
	// 		Age  int    `json:"age"`
	// 	}

	// 	u := user{}
	// 	err := decoder.Decode(&u)
	// 	if err != nil {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body"})
	// 		return
	// 	}

	// 	err = validator.New().Struct(u)
	// 	if err != nil {
	// 		w.Header().Set("Content-Type", "application/json")
	// 		w.WriteHeader(http.StatusBadRequest)
	// 		_ = json.NewEncoder(w).Encode(map[string]string{"message": "invalid request body"})
	// 		return
	// 	}

	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusOK)

	// 	_ = json.NewEncoder(w).Encode(map[string]string{"message": u.Name})
	// })

	router.GET("/inventories", inventoryCtrl.GetAll)
	router.GET("/inventories/:code", inventoryCtrl.GetByCode)
	router.POST("/inventories", inventoryCtrl.Create)
	router.PUT("/inventories/:code", inventoryCtrl.Update)
	router.DELETE("/inventories/:code", inventoryCtrl.Delete)

	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`("data":{}, "message":"route not found")`))
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"data":{}, "message":"method not allowed"}`))
	})

	router.PanicHandler = func(w http.ResponseWriter, r *http.Request, err interface{}) {
		logger.Error("Panic Handler", slog.Any("error", err))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{"message": http.StatusText(http.StatusInternalServerError)})
	}

	logger.Info("API service running in " + config.AppHost + ":" + config.AppPort)
	server := &http.Server{
		Addr:    config.AppHost + ":" + config.AppPort,
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
