package database

import (
	"errors"
	"testing"

	"gorm.io/gorm"
)

func TestAddUserToServer(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    userID, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    name := "guld"
    serverID, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting new server", err)

    err = db.AddUserToServer(9999, serverID, ROLE_SERVER_ADMIN)
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        t.Fatal("Should not be able to add to non-existent user")
    }
    err = db.AddUserToServer(userID, 9999, ROLE_SERVER_ADMIN)
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        t.Fatal("Should not be able to add to non-existent server")
    }

    err = db.AddUserToServer(userID, serverID, ROLE_SERVER_ADMIN)
    checkErr(t, "Error adding user to server", err)

    users, err := db.GetAllUsersByServerID(serverID)
    checkErr(t, "Error getting users by server", err)
    if len(users) != 1 {
        t.Fatal("Unexpected number of users")
    }
    if users[0].Username != username || users[0].Password != password {
        t.Fatal("Faulty user info")
    }

    userServer, err := db.GetUserServerByIDs(userID, serverID)
    checkErr(t, "Error getting UserServer", err)

    if userServer.Role != ROLE_SERVER_ADMIN {
        t.Fatal("UserServer has the wrong role")
    }

    name2 := "silver"
    serverID2, err := db.InsertNewServer(name2)
    checkErr(t, "Error inserting new server2", err)
    err = db.AddUserToServer(userID, serverID2, ROLE_NORMAL)
    checkErr(t, "Error adding users to server2", err)

    username2 := "jack"
    password2 := "pass"
    userID2, err := db.InsertNewUser(username2, password2)
    checkErr(t, "Error inserting user2", err)
    
    err = db.AddUserToServer(userID2, serverID2, ROLE_SERVER_ADMIN)
    checkErr(t, "Error adding user2 to server2", err)

    users, err = db.GetAllUsersByServerID(serverID2)
    checkErr(t, "Error getting users by server2", err)

    if len(users) != 2 {
        t.Fatalf("Unexpected number of users. Users: %+v\n", users)
    }

    servers, err := db.GetAllServersByUserID(userID)
    checkErr(t, "Error getting servers by user", err)

    if len(servers) != 2 {
        t.Fatal("Unexpected number of servers")
    }
}

func TestRemoveUserFromServer(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    userID, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    name := "guld"
    serverID, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting new server", err)

    err = db.AddUserToServer(userID, serverID, ROLE_SERVER_ADMIN)
    checkErr(t, "Error adding user to server", err)

    servers, err := db.GetAllServersByUserID(userID)
    checkErr(t, "Error getting servers by userID", err)

    if len(servers) != 1 {
        t.Fatalf("User belong to an unexpected number of servers. %+v\n", servers)
    }
    
    err = db.RemoveUserFromServer(userID, serverID)
    checkErr(t, "Error removing user from server", err)

    servers, err = db.GetAllServersByUserID(userID)
    checkErr(t, "Error getting servers by userID", err)

    if len(servers) != 0 {
        t.Fatal("User belong to an unexpected number of servers")
    }
}

func TestChangeUserRoleInServer(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    userID, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    name := "guld"
    serverID, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting new server", err)

    err = db.AddUserToServer(userID, serverID, ROLE_SERVER_ADMIN)
    checkErr(t, "Error adding user to server", err)

    userServer, err := db.GetUserServerByIDs(userID, serverID)
    checkErr(t, "Error getting UserServer", err)

    if userServer.Role != ROLE_SERVER_ADMIN {
        t.Fatal("UserServer has the wrong role")
    }
    
    err = db.ChangeUserRoleInServer(userID, serverID, ROLE_NORMAL)
    checkErr(t, "Error changing role", err)


    userServer, err = db.GetUserServerByIDs(userID, serverID)
    checkErr(t, "Error getting UserServer", err)

    if userServer.Role != ROLE_NORMAL {
        t.Fatal("UserServer has the wrong role")
    }
}
