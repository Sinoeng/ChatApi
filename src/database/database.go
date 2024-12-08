package database

import (
	"fmt"
	"log"
	"os"

	mysqlErr "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ChatApiDB struct {
	db *gorm.DB // connection pool
}

func InitDatabase() (*ChatApiDB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"),
	)
	gormConfig := gorm.Config{PrepareStmt: true, Logger: logger.Discard}
	// opens the connection to the database. second arg is configurations
	db, err := gorm.Open(mysql.Open(dsn), &gormConfig) // db is a connection pool
    switch err.(type){
    case *mysqlErr.MySQLError:
        db_name := os.Getenv("MYSQL_DATABASE")

        query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", db_name)
        dsn2 := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
            os.Getenv("MYSQL_USER"),
            os.Getenv("MYSQL_PASSWORD"),
            os.Getenv("DB_HOST"),
            os.Getenv("MYSQL_PORT"),
        )

        db, err = gorm.Open(mysql.Open(dsn2), &gormConfig)
        if err != nil {
            log.Fatalf("Establishing connection to %s failed due to %s\n", db_name, err.Error())
            break
        }
        err = db.Exec(query).Error
        if err != nil {
            log.Fatalf("Creating database %s failed due to %s\n", db_name, err.Error())
            break
        }
        mdb, _ := db.DB()
        mdb.Close()
        db, err = gorm.Open(mysql.Open(dsn), &gormConfig)
        if err != nil {
            return nil, err
        }
        break
    case nil:
        break
    default:
        log.Printf("Error establishing connection: %s\n", err.Error())
        return nil, err
    }

	chatApiDB := ChatApiDB{
		db: db,
	}
	err = chatApiDB.autoMigration()
	if err != nil {
        log.Printf("Error in automigration. Error: %s\n", err.Error())
		return &chatApiDB, err
	}

	return &chatApiDB, err
}

func (self *ChatApiDB) autoMigration() error {
	return self.db.AutoMigrate(
        &Message{},
        &Server{},
		&User{},
        &UserServer{},
	)
}
