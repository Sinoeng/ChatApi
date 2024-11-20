package database
//
// import (
//     "gorm.io/driver/mysql"
//     "gorm.io/gorm"
//     "fmt"
//     "os"
// )
//
// type ChatApiDB struct {
//     db *gorm.DB
// }
//
//
// func InitDatabase() (*ChatApiDB, error) {
// 	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
// 		os.Getenv("MYSQL_USER"),
// 		os.Getenv("MYSQL_PASSWORD"),
// 		os.Getenv("DB_HOST"),
// 		os.Getenv("MYSQL_PORT"),
// 		os.Getenv("MYSQL_DATABASE"),
// 	)
// 	gormConfig := gorm.Config{PrepareStmt: true}
// 	// opens the connection to the database. second arg is configurations
// 	db, err := gorm.Open(mysql.Open(dsn), &gormConfig)
//     if err != nil {
//         return nil, err
//     }
//
//     chatApiDB := ChatApiDB{
//         db: db,
//     }
//     err = chatApiDB.autoMigration()
//     if err != nil {
//         return &chatApiDB, err
//     }
//
//     return &chatApiDB, err
// }
//
// func (self *ChatApiDB) autoMigration() error {
//     return self.db.AutoMigrate()
// }
//
// func (self *ChatApiDB) Getter(){
//     self.db.Exec("SHOW DATABASES;")
// }
