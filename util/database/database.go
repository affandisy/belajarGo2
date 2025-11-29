package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	DBDriver      string
	DBEnableDebug bool

	DBMySQLHost     string
	DBMySQLPort     string
	DBMySQLUser     string
	DBMySQLPassword string
	DBMySQLName     string

	DBPostgreSQLHost     string
	DBPostgreSQLPort     string
	DBPostgreSQLUser     string
	DBPostgreSQLPassword string
	DBPostgreSQLName     string

	DBSQLiteName string
}

func (conf *Config) GetDatabaseConnection() *gorm.DB {
	var err error
	var db *gorm.DB
	// conf.DBEnableDebug = true

	switch conf.DBDriver {
	case "mysql":
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			conf.DBMySQLUser,
			conf.DBMySQLPassword,
			conf.DBMySQLHost,
			conf.DBMySQLPort,
			conf.DBMySQLName,
		)

		db, err = gorm.Open(mysql.Open(dsn))
		if err != nil {
			log.Fatal(err)
		}
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(conf.DBSQLiteName))
		if err != nil {
			log.Fatal(err)
		}
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=require",
			conf.DBPostgreSQLHost,
			conf.DBPostgreSQLPort,
			conf.DBPostgreSQLUser,
			conf.DBPostgreSQLPassword,
			conf.DBPostgreSQLName,
		)
		db, err = gorm.Open(postgres.Open(dsn))
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("unsupported driver")
	}

	if err != nil {
		log.Fatal(err)
	}

	if conf.DBEnableDebug {
		return db.Debug()
	}

	return db
}
