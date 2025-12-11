package main

import (
	invHandler "belajarGo2/app/echo-server/controller/inventory"
	userController "belajarGo2/app/echo-server/controller/user"
	"belajarGo2/app/echo-server/router"
	invRepo "belajarGo2/repository/inventory"
	"belajarGo2/repository/notification/mailjet"
	userRepo "belajarGo2/repository/user"
	invSvc "belajarGo2/service/inventory"
	userService "belajarGo2/service/user"
	"belajarGo2/util/database"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	cfg "github.com/pobyzaarif/go-config"
	echoSwagger "github.com/swaggo/echo-swagger"
)

var loggerOption = slog.HandlerOptions{AddSource: true}
var logger = slog.New(slog.NewJSONHandler(os.Stdout, &loggerOption))

type Config struct {
	// AppHost                 string `env:"APP_HOST"`
	// AppPort                 string `env:"APP_PORT"`
	AppHost                 string `env:"APP_PHOST"`
	AppPort                 string `env:"APP_PORT_ECHO_SERVER"`
	AppDeploymentUrl        string `env:"APP_DEPLOYMENT_URL"`
	AppEmailVerificationKey string `env:"APP_EMAIL_VERIFICATION_KEY"`
	AppJWTSecret            string `env:"APP_JWT_SECRET"`

	DBDriver        string `env:"DB_DRIVER"`
	DBMySQLHost     string `env:"DB_MYSQL_HOST"`
	DBMySQLPort     string `env:"DB_MYSQL_PORT"`
	DBMySQLUser     string `env:"DB_MYSQL_USER"`
	DBMySQLPassword string `env:"DB_MYSQL_PASSWORD"`
	DBMySQLName     string `env:"DB_MYSQL_NAME"`

	DBSQLiteName string `env:"DB_SQLITE_NAME"`

	DBPostgreSQLHost     string `env:"DB_POSTGRESQL_HOST"`
	DBPostgreSQLPort     string `env:"DB_POSTGRESQL_PORT"`
	DBPostgreSQLUser     string `env:"DB_POSTGRESQL_USER"`
	DBPostgreSQLPassword string `env:"DB_POSTGRESQL_PASSWORD"`
	DBPostgreSQLName     string `env:"DB_POSTGRESQL_NAME"`

	DBMongoURI  string `env:"DB_MONGO_URI"`
	DBMongoName string `env:"DB_MONGO_NAME"`

	MailjetBaseUrl           string `env:"MAILJET_BASE_URL"`
	MailjetBasicAuthUsername string `env:"MAILJET_BASIC_AUTH_USERNAME"`
	MailjetBasicAuthPassword string `env:"MAILJET_BASIC_AUTH_PASSWORD"`
	MailjetSenderEmail       string `env:"MAILJET_SENDER_EMAIL"`
	MailjetSenderName        string `env:"MAILJET_SENDER_NAME"`
}

func main() {
	spew.Dump()

	config := Config{}
	cfg.LoadConfig(&config)
	logger.Info("Config loaded")

	// Init db connection
	databaseConfig := database.Config{
		DBDriver:             config.DBDriver,
		DBMySQLHost:          config.DBMySQLHost,
		DBMySQLPort:          config.DBMySQLPort,
		DBMySQLUser:          config.DBMySQLUser,
		DBMySQLPassword:      config.DBMySQLPassword,
		DBMySQLName:          config.DBMySQLName,
		DBSQLiteName:         config.DBSQLiteName,
		DBPostgreSQLHost:     config.DBPostgreSQLHost,
		DBPostgreSQLPort:     config.DBPostgreSQLPort,
		DBPostgreSQLUser:     config.DBPostgreSQLUser,
		DBPostgreSQLPassword: config.DBPostgreSQLPassword,
		DBPostgreSQLName:     config.DBPostgreSQLName,
		DBMongoURI:           config.DBMongoURI,
		DBMongoName:          config.DBMongoName,
	}

	db := databaseConfig.GetDatabaseConnection()
	logger.Info("Database client connected!")

	// Setup server
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.CORS())
	e.Use(middleware.LoggerWithConfig(
		middleware.LoggerConfig{
			Skipper: middleware.DefaultSkipper,
			Format: `{"time":"${time_rfc3339_nano}","level":"INFO","id":"${id}","remote_ip":"${remote_ip}",` +
				`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
				`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
				`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
			CustomTimeFormat: "2006-01-02 15:04:05.00000",
		},
	))
	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(middleware.Recover())

	// Setup routes
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": "pong",
		})
	})

	// Swagger
	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())

	// notification
	mailjetEmail := mailjet.NewMailjetRepository(
		logger,
		mailjet.MailjetConfig{
			MailjetBaseURL:           config.MailjetBaseUrl,
			MailjetBasicAuthUsername: config.MailjetBasicAuthUsername,
			MailjetBasicAuthPassword: config.MailjetBasicAuthPassword,
			MailjetSenderEmail:       config.MailjetSenderEmail,
			MailjetSenderName:        config.MailjetSenderName,
		},
	)

	// inventory endpoint
	inventoryRepo := invRepo.NewGormRepository(db)
	inventorySvc := invSvc.NewService(inventoryRepo)
	inventoryCtrl := invHandler.NewController(logger, inventorySvc)

	// endpoint
	// e.GET("/inventory", inventoryCtrl.GetAll)
	// e.GET("/inventories/:code", inventoryCtrl.GetByCode)
	// e.POST("/inventories", inventoryCtrl.Create)
	// e.PUT("/inventories/:code", inventoryCtrl.Update)
	// e.DELETE("/inventories/:code", inventoryCtrl.Delete)

	// endpoint group inventory
	// inventoryEndpoint := e.Group("/inventories")
	// inventoryEndpoint.GET("", inventoryCtrl.GetAll)
	// inventoryEndpoint.GET("/:code", inventoryCtrl.GetByCode)
	// inventoryEndpoint.POST("", inventoryCtrl.Create)
	// inventoryEndpoint.PUT("/:code", inventoryCtrl.Update)
	// inventoryEndpoint.DELETE("/:code", inventoryCtrl.Delete)

	// endpoint user
	// userRepo := userRepo.NewGormRepository(db)

	// using mongodb
	dbMongo := databaseConfig.GetNoSQLDatabaseConnection()
	userMongoRepo := userRepo.NewMongoRepository(dbMongo)

	userService := userService.NewService(logger, userMongoRepo, config.AppDeploymentUrl, config.AppJWTSecret, config.AppEmailVerificationKey, mailjetEmail)
	userCtrl := userController.NewController(logger, userService)

	// endpoint group user
	// userEndpoint := e.Group("/users")
	// userEndpoint.POST("/register", userCtrl.Register)
	// userEndpoint.POST("/login", userCtrl.Login)

	router.RegisterPath(e, config.AppJWTSecret, inventoryCtrl, userCtrl)

	// Start server
	address := config.AppHost + ":" + config.AppPort
	go func() {
		if err := e.Start(address); err != http.ErrServerClosed {
			log.Fatal("Failed on http server " + config.AppPort)
		}
	}()

	logger.Info("API service running in " + address)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	// a timeout of 10 seconds to shutdown the server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(ctx); err != nil {
		log.Fatal("Failed to shutting down echo server", "error", err)
	} else {
		logger.Info("Successfully shutting down echo server")
	}
}
