package config

import (
	"fmt"
	"log"
	"os"

	"shedoo-backend/internal/models"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

var dbInstance *Service

type RedisConfig struct {
	Addr     string
	Password string
}

func NewRedis() *RedisConfig {
	return &RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
	}
}

func NewDB() *Service {
	if dbInstance != nil {
		return dbInstance
	}

	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USERNAME")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DATABASE")
	adminAccount := os.Getenv("ADMIN_ACCOUNT")

	dsn := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}
	err = db.AutoMigrate(
		&models.Course{},
		&models.Enrollment{},
		&models.ScrapeCourseJob{},
		&models.ScrapeExamJob{},
		&models.CourseExam{},
		&models.Admin{},
	)
	if err != nil {
		panic("failed to auto migrate models")
	}

	if adminAccount != "" {
		admin := models.Admin{Account: adminAccount}
		if err := db.Where("account = ?", adminAccount).FirstOrCreate(&admin).Error; err != nil {
			log.Printf("Failed to initialize admin (%s): %v", adminAccount, err)
		} else {
			log.Printf("Admin initialized or already exists: %s", adminAccount)
		}
	} else {
		log.Println("ADMIN_ACCOUNT not set, skipping admin initialization")
	}

	dbInstance = &Service{DB: db}
	return dbInstance
}
