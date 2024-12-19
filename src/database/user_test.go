package database

import (
	"errors"
	"testing"

	mysqlErr "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func TestInsertNewUser(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    _, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    username2 := "john2"
    password2 := "pass"
    _, err = db.InsertNewUser(username2, password2)
    checkErr(t, "Error inserting user2", err)


    _, err = db.InsertNewUser(username2, password2)
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
    users, err := db.GetAllUsers()
    checkErr(t, "Error getting users", err)
    if len(users) != 2 {
        t.Fatal("Unexpected number of users")
    }
    if users[0].Username != username || users[0].Password != password || users[0].Email != "" {
        t.Fatalf("User1 data is incorrect %+v", users[0])
    }
    if users[1].Username != username2 || users[1].Password != password2 || users[1].Email != "" {
        t.Fatalf("User2 data is incorrect %+v", users[1])
    }
}

func TestInsertNerUserWEmail(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    email := "john.son@gmail.com"
    _, err := db.InsertNewUserWEmail(username, password, email)
    checkErr(t, "Error inserting user", err)

    username2 := "john2"
    password2 := "pass"
    email2 := "john.son2@gmail.com"
    _, err = db.InsertNewUserWEmail(username2, password2, email2)
    checkErr(t, "Error inserting user2", err)

    _, err = db.InsertNewUser(username2, password2)
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
    users, err := db.GetAllUsers()
    checkErr(t, "Error getting users", err)
    if len(users) != 2 {
        t.Fatal("Unexpected number of users")
    }
    if users[0].Username != username || users[0].Password != password || users[0].Email != email {
        t.Fatalf("User1 data is incorrect %+v\n", users[0])
    }
    if users[1].Username != username2 || users[1].Password != password2 || users[1].Email != email2 {
        t.Fatalf("User2 data is incorrect %+v\n", users[1])
    }
}

func TestChangeEmailByID(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    email := "john.son@gmail.com"
    email2 := "john.son2@gmail.com"
    _, err := db.InsertNewUserWEmail(username, password, email)
    checkErr(t, "Error inserting user", err)

    user, err := db.GetUserByUsername(username)
    checkErr(t, "Error getting user", err)
    err = db.ChangeUserEmailByID(user.ID, email2)

    user, err = db.GetUserByUsername(username)
    checkErr(t, "Error getting user2", err)
    if user.Email != email2 {
        t.Fatal("Failed to change email")
    }
}

func TestDeleteUserByID(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    email := "john.son@gmail.com"
    userID, err := db.InsertNewUserWEmail(username, password, email)
    checkErr(t, "Error inserting user", err)

    name := "guld"
    serverID, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting new server", err)

    err = db.AddUserToServer(userID, serverID, ROLE_NORMAL)
    checkErr(t, "Error adding user to server", err)

    users, err := db.GetAllUsers()
    checkErr(t, "Error getting users", err)
    if len(users) != 1 {
        t.Fatal("Unexpected number of users")
    }
    err = db.DeleteUserByID(users[0].ID)
    checkErr(t, "Error deleting user", err)

    users, err = db.GetAllUsers()
    checkErr(t, "Error getting users", err)
    if len(users) != 0 {
        t.Fatal("User was not deleted")
    }

    users, err = db.GetAllUsersByServerID(serverID)
    checkErr(t, "Error getting users by server", err)

    if len(users) != 0 {
        t.Fatal("User is still in server")
    }
    userServer, err := db.GetUserServerByIDs(userID, serverID)
    if !errors.Is(err, gorm.ErrRecordNotFound) {
        t.Fatalf("UserServer not deleted. userServer: %+v\n", userServer)
    }

    server, err := db.GetServerByID(serverID)
    checkErr(t, "Error getting server", err)

    if len(server.Users) != 0 {
        t.Fatal("Server still holds user")
    }
}

func TestGetUserByID(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    email := "john.son@gmail.com"
    _, err := db.InsertNewUserWEmail(username, password, email)
    checkErr(t, "Error inserting user", err)

    users, err := db.GetAllUsers()
    checkErr(t, "Error getting users", err)
    if len(users) != 1 {
        t.Fatal("Unexpected number of users")
    }
    user, err := db.GetUserByID(users[0].ID)
    checkErr(t, "Error getting user", err)

    if user.ID != users[0].ID || user.Username != users[0].Username || user.Password != users[0].Password || user.Email != users[0].Email {
        t.Fatal("User was not identical")
    }
}
