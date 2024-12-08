package database

import (
	"testing"

	mysqlErr "github.com/go-sql-driver/mysql"
)

func TestInsertNewUser(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    username2 := "john2"
    password2 := "pass"
    err = db.InsertNewUser(username2, password2)
    checkErr(t, "Error inserting user2", err)


    err = db.InsertNewUser(username2, password2)
    switch e := err.(type) {
    case *mysqlErr.MySQLError:
        switch e.Number {
        case 1062:
            // OK
            break
        default:
            t.Fatalf("Unrecognised error number. Error: %s\n", err.Error())
        }
    default:
        t.Fatalf("Unrecognised error. Error: %s\n", err.Error())
    }
}
