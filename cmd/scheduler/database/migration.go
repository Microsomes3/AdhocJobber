package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetDatabaseConnection() (*gorm.DB, error) {

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/scheduler", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_IP"), os.Getenv("DB_PORT"))

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
