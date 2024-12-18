package database

import (
	"errors"
	"testing"
	"gorm.io/gorm"
)

func TestInsertNewServer(t *testing.T) {
    resetDB()
    name := "guld"
    id, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting new server", err)
    
    server, err := db.GetServerByID(id)
    checkErr(t, "Error getting server", err)
    if server.Name != name {
        t.Fatal("Server names do not match")
    }
}

func TestDeleteServer(t *testing.T) {
    resetDB()
    name := "guld"
    id, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting new server", err)
    
    server, err := db.GetServerByID(id)
    checkErr(t, "Error getting server", err)
    if server.Name != name {
        t.Fatal("Server names do not match")
    }

    err = db.DeleteServerByID(id)
    checkErr(t, "Error deleting server", err)

    server, err = db.GetServerByID(id)
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        t.Fatalf("Incorrect error received. Error: %s\n", err)
    }
}

func TestChangeServerName(t *testing.T) {
    resetDB()
    name := "guld"
    id, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting new server", err)

    name2 := "silver"
    err = db.ChangeServerName(id, name2)
    checkErr(t, "Error changeing server name", err)

    server, err := db.GetServerByID(id)
    checkErr(t, "Error getting server", err)

    if server.Name != name2 {
        t.Fatal("Server name was not changed")
    }
}
