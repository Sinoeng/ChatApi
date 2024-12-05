package database

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UserDB struct {
	db *gorm.DB // connection pool
}

func InitDatabase() (*UserDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	gormConfig := gorm.Config{PrepareStmt: true}
	// opens the connection to the database. second arg is configurations
	db, err := gorm.Open(mysql.Open(dsn), &gormConfig) // db is a connection pool
	if err != nil {
		return nil, err
	}

	chatApiDB := UserDB{
		db: db,
	}
	err = chatApiDB.autoMigration()
	if err != nil {
		return &chatApiDB, err
	}

	return &chatApiDB, err
}

func (self *UserDB) autoMigration() error {
	return self.db.AutoMigrate(
		&User{},
	)
}
