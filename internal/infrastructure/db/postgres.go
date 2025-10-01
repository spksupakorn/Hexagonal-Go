package database

import (
	"dungeons-dragon-service/internal/config"
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	dbInstance *postgresDatabase
)

func NewPostgresDatabase() Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			config.GetConfigString("DB_HOST"),
			config.GetConfigString("DB_USER"),
			config.GetConfigString("DB_PASSWORD"),
			config.GetConfigString("DB_NAME"),
			config.GetConfigInt("DB_PORT"),
			config.GetConfigString("DB_SSLMODE"),
			config.GetConfigString("DB_TIMEZONE"),
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		//log info about successful connection
		fmt.Println("⚔️  Successfully connected to the database ⚔️ ")

		dbInstance = &postgresDatabase{Db: db}
	})

	return dbInstance
}

func (p *postgresDatabase) ConnectDB() *gorm.DB {
	return p.Db
}
