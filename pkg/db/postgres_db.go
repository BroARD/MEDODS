package db

import (
	"Medods/internal/models"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitDB() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("FATAL: Ошибка загрузки .env файла: %v", err)
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	sslmode := os.Getenv("SSL_MODE")
	
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("FATAL: Ошибка подключения к БД: %v", err)
	}

	//TODO: Добавить миграции
	if err := db.AutoMigrate(&models.RefreshToken{}); err != nil {
		log.Fatalf("FATAL: Ошибка миграции: %v", err)
	}

	return db, nil
}
