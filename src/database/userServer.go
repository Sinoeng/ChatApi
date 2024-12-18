package database

type UserServer struct {
    UserID uint `gorm:"primaryKey"`
    ServerID uint `gorm:"primaryKey"`
    Role string
}

func (self *ChatApiDB) AddUserToServer(userID uint, serverID uint, role string) error {
    userServer := UserServer{
        UserID: userID,
        ServerID: serverID,
        Role: role,
    }
    return self.db.Create(&userServer).Error
}

func (self *ChatApiDB) GetAllUsersByServerID(serverID uint) ([]User, error) {
    var res []User
    var tmp []uint
    err := self.db.Model(&UserServer{}).Select("user_id").Where(&UserServer{ServerID: serverID}).Find(&tmp).Error
    if err != nil {
        return res, err
    }
    err = self.db.Where("id IN ?", tmp).Find(&res).Error
    return res, err
    // var server Server
    // // err := self.db.Preload("Users").Where(&Server{ID: serverID}).First(&server).Error
    // err := self.db.Preload(clause.Associations).Where(&Server{ID: serverID}).First(&server).Error
    // log.Printf("Server: %+v", server)
    // var tst UserServer
    // self.db.Where(&UserServer{ServerID: serverID}).Find(&tst)
    // log.Printf("tst: %+v", tst)
    // return server.Users, err
}

func (self *ChatApiDB) GetAllServersByUserID(userID uint) ([]Server, error) {
    // var user User
    // err := self.db.Preload("Servers").Where(&User{ID: userID}).First(&user).Error
    // log.Printf("User: %+v", user)
    // return user.Servers, err
    var res []Server
    var tmp []uint
    err := self.db.Model(&UserServer{}).Select("server_id").Where(&UserServer{UserID: userID}).Find(&tmp).Error
    if err != nil {
        return res, err
    }
    err = self.db.Where("id IN ?", tmp).Find(&res).Error
    return res, err
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

