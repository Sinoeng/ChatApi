package database

import (
	"os"
	"testing"
)

var db *ChatApiDB

func TestMain(m *testing.M) {
    initDB()

    m.Run()

    db.ResetTestDBDB()
}

func initDB() {
    var err error
    db, err = InitDatabase(os.Getenv("MYSQL_DATABASE")+"DB")
    if err != nil {
        println("Panicing in init DB")
        panic(err)
    }
}

func resetDB() {
    db.ResetTestDBDB()
    initDB()
}

func checkErr(t *testing.T, desc string, err error) {
    if err != nil {
        t.Fatalf(desc + ". Err: %s\n", err.Error())
    }
}
