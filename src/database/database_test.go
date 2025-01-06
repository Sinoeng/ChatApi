package database

import "testing"

var db *ChatApiDB

func TestMain(m *testing.M) {
    initDB()

    m.Run()

    dropDB()
}

func initDB() {
    var err error
    db, err = InitDatabase()
    if err != nil {
        println("Panicing in init DB")
        panic(err)
    }
}

func dropDB() {
    db.db.Exec("DROP DATABASE ChatApiTest;")
}

func resetDB() {
    dropDB()
    initDB()
}

func checkErr(t *testing.T, desc string, err error) {
    if err != nil {
        t.Fatalf(desc + ". Err: %s\n", err.Error())
    }
}
