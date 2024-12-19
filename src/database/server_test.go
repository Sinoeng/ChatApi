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

    username := "john"
    password := "pass"
    userID, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting new user", err)

    err = db.AddUserToServer(userID, id, ROLE_NORMAL)
    checkErr(t, "Error adding user to server", err)
    
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

    _, err = db.GetUserServerByIDs(userID, id)
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        t.Fatal("UserServer still exists")
    }

    user, err := db.GetUserByID(userID)
    checkErr(t, "Error getting user by id", err)

    if len(user.Servers) != 0 {
        t.Fatal("User still holds the server")
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
