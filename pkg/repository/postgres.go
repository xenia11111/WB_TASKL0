package repository

import (
	"fmt"

	"github.com/sirupsen/logrus"
	app "github.com/xenia11111/WB_TASKL0"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	ordersTable = "orders"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		logrus.Fatal(err)
	}

	db.AutoMigrate(&app.Order{})

	return db, nil
}
