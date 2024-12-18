package database

import "testing"

func TestAddUserToServer(t *testing.T) {
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

    name2 := "silver"
    serverID, err = db.InsertNewServer(name2)
    checkErr(t, "Error inserting new server2", err)
    err = db.AddUserToServer(userID, serverID, ROLE_NORMAL)
    checkErr(t, "Error adding users to server2", err)

    username2 := "jack"
    password2 := "pass"
    userID2, err := db.InsertNewUser(username2, password2)
    checkErr(t, "Error inserting user", err)
    
    err = db.AddUserToServer(userID2, serverID, ROLE_SERVER_ADMIN)
    checkErr(t, "Error adding user2 to server", err)

    users, err := db.GetAllUsersByServerID(serverID)
    checkErr(t, "Error getting users by server", err)

    if len(users) != 2 {
        t.Fatal("Unexpected number of users")
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

    // user, err := db.GetUserByID(userID)
    // checkErr(t, "Error getting users by id", err)
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
