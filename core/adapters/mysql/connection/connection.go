package connection

import (
	"biblebrain-domain-go/core/domain"

	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetConnection() (*gorm.DB, error) {
	cfg, err := domain.LoadConfig()
	if err != nil {
		fmt.Println("Error load database config:", err)
		return nil, err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return nil, err
	}
	fmt.Println("Connected to the database.")
	return db, nil
}
