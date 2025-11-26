package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBDriver        string
	DBMySQLHost     string
	DBMySQLPort     string
	DBMySQLUser     string
	DBMySQLPassword string
	DBMySQLName     string

	DBSQLiteName string
}

func (conf *Config) GetDatabaseConnection() *gorm.DB {
	if conf.DBDriver == "mysql" {
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.DBMySQLUser,
			conf.DBMySQLPassword,
			conf.DBMySQLHost,
			conf.DBMySQLPort,
			conf.DBMySQLName,
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newDBLogger()})
		if err != nil {
			log.Fatal(err)
		}

		return db.Debug()
	}

	if conf.DBDriver == "sqlite" {
		db, err := gorm.Open(sqlite.Open(conf.DBSQLiteName), &gorm.Config{Logger: newDBLogger()})
		if err != nil {
			log.Fatal(err)
		}

		return db.Debug()
	}

	log.Fatal("unsupported driver")

	return nil
}

func newDBLogger() logger.Interface {
	return logger.New(
		log.Default(),
		logger.Config{
			SlowThreshold:             30 * time.Second, // Slow SQL threshold
			LogLevel:                  logger.Silent,    // Log level
			IgnoreRecordNotFoundError: false,            // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,            // Enable Color
		},
	)
}
