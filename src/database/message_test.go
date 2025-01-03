package database

import (
	"testing"

	mysqlErr "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

func TestNewMessage(t *testing.T) {
    resetDB()
    payload := "Eine messag"
    username := "john"
    password := "pass"
    userID, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    name := "guld"
    serverID, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting server", err)

    checkForeignKey := func (t *testing.T, err error)  {
        switch e := err.(type) {
        case *mysqlErr.MySQLError:
            switch e.Number {
            case 1452:
                //foreign key constraint fail
                break
            default:
                t.Fatalf("Unknown mysql error. Error %s\n", e.Error())
            }
        default:
            t.Fatalf("Unknown error. Error %s\n", e.Error())
        }
    }
    _, err = db.NewMessage(payload, 9999, 9999)
    checkForeignKey(t, err)

    _, err = db.NewMessage(payload, userID, 9999)
    checkForeignKey(t, err)

    _, err = db.NewMessage(payload, 9999, serverID)
    checkForeignKey(t, err)

    msgID, err := db.NewMessage(payload, userID, serverID)
    checkErr(t, "Error creating new message", err)

    msg, err := db.GetMessageByID(msgID)
    checkErr(t, "Error getting message by id", err)

    if msg.Payload != payload {
        t.Fatal("Message payload is faulty")
    }

    msgs, err := db.GetMessagesByUserID(userID)
    checkErr(t, "Error getting messages by userID", err)

    if len(msgs) != 1 {
        t.Fatal("Unexpected number of messages")
    }

    if msgs[0].Payload != payload {
        t.Fatal("Faulty message payload")
    }

    msgs, err = db.GetMessagesByServerID(serverID)
    checkErr(t, "Error getting messages by serverID", err)

    if len(msgs) != 1 {
        t.Fatal("Unexpected number of messages")
    }

    if msgs[0].Payload != payload {
        t.Fatal("Faulty message payload")
    }
}

func TestDeleteMessage(t *testing.T) {
    resetDB()
    payload := "Eine messag"
    username := "john"
    password := "pass"
    userID, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    name := "guld"
    serverID, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting server", err)

    msgID, err := db.NewMessage(payload, userID, serverID)
    checkErr(t, "Error creating new message", err)

    _, err = db.GetMessageByID(msgID)
    checkErr(t, "Error getting message by id", err)

    err = db.DeleteMessage(msgID)
    checkErr(t, "Error deleting message", err)

    _, err = db.GetMessageByID(msgID)
    switch err {
    case gorm.ErrRecordNotFound:
        break
    case nil:
        t.Fatal("No error received")
    default:
        t.Fatalf("Faulty error. Error: %s\n", err.Error())
    }
}

func TestGetMessagesByUserIDAndServerID(t *testing.T) {
    resetDB()
    username := "john"
    password := "pass"
    userID, err := db.InsertNewUser(username, password)
    checkErr(t, "Error inserting user", err)

    username2 := "jack"
    password2 := "pass"
    userID2, err := db.InsertNewUser(username2, password2)
    checkErr(t, "Error inserting user", err)

    name := "guld"
    serverID, err := db.InsertNewServer(name)
    checkErr(t, "Error inserting server", err)

    err = db.AddUserToServer(userID, serverID, ROLE_NORMAL)
    checkErr(t, "Error adding user1 to server1", err)
    err = db.AddUserToServer(userID2, serverID, ROLE_NORMAL)
    checkErr(t, "Error adding user1 to server1", err)

    payload := "content"
    _, err = db.NewMessage(payload, userID, serverID)
    checkErr(t, "Error leaving message", err)
    _, err = db.NewMessage(payload, userID, serverID)
    checkErr(t, "Error leaving message", err)
    _, err = db.NewMessage(payload, userID, serverID)
    checkErr(t, "Error leaving message", err)
    _, err = db.NewMessage(payload, userID2, serverID)
    checkErr(t, "Error leaving message", err)

    msgs, err := db.GetMessagesByUserIDAndServerID(userID, serverID)
    checkErr(t, "Error getting messages by userID and serverID", err)

    if len(msgs) != 3 {
        t.Fatal("Unexpected number of messages")
    }

    msgs, err = db.GetMessagesByUserIDAndServerID(userID2, serverID)
    checkErr(t, "Error getting messages by userID2 and serverID", err)

    if len(msgs) != 1 {
        t.Fatal("Unexpected number of messages")
    }
}
