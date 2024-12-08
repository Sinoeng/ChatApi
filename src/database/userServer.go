package database

type UserServer struct {
    UserID uint `gorm:"primaryKey"`
    ServerID uint `gorm:"primaryKey"`
    Role string //change?
}

func (self *ChatApiDB) AddUserToServer(userID uint, serverID uint, role string) error {
    userServer := UserServer{
        UserID: userID,
        ServerID: serverID,
        Role: role,
    }
    return self.db.Create(&userServer).Error
}

func (self *ChatApiDB) RemoveUserFromServer(userID uint, serverID uint) error {
    return self.db.Where(&UserServer{UserID: userID, ServerID: serverID}).Delete(UserServer{}).Error
}

func (self *ChatApiDB) ChangeUserRoleInServer(userID uint, serverID uint, newRole string) error {
    tx := self.db.Begin()
    var userServer UserServer
    err := tx.Where(&UserServer{UserID: userID, ServerID: serverID}).First(&userServer).Error
    if err != nil {
        tx.Rollback()
        return err
    }
    userServer.Role = newRole
    err = tx.Save(userServer).Error
    if err != nil {
        tx.Rollback()
        return err
    }
    err = tx.Commit().Error
    if err != nil {
        tx.Rollback()
    }
    return err
}

