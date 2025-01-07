package database

import "gorm.io/gorm"

type UserServer struct {
	UserID   uint `gorm:"primaryKey"`
	ServerID uint `gorm:"primaryKey"`
	Role     string
}

func (self *ChatApiDB) AddUserToServer(userID uint, serverID uint, role string) error {
	// var user User
	// var server Server
	// tx := self.db.Begin()
	// err := tx.First(&user, userID).Error
	// if err != nil {
	//     tx.Rollback()
	//     return err
	// }
	// err = tx.First(&server, serverID).Error
	// if err != nil {
	//     tx.Rollback()
	//     return err
	// }
	// userServer := UserServer{
	//     UserID: userID,
	//     ServerID: serverID,
	//     Role: role,
	// }
	// err = tx.Create(&userServer).Error
	// if err != nil {
	//     tx.Rollback()
	//     return err
	// }
	// err = tx.Commit().Error
	// return err

	var user User
	var server Server
	tx := self.db.Begin()
	err := tx.First(&user, userID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.First(&server, serverID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&server).Association("Users").Append(&user)
	if err != nil {
		tx.Rollback()
		return err
	}
	var userServer UserServer
	err = tx.Where(&UserServer{UserID: userID, ServerID: serverID}).First(&userServer).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	userServer.Role = role
	err = tx.Save(&userServer).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return err

}

func (self *ChatApiDB) GetUserServerByIDs(userID, serverID uint) (UserServer, error) {
	var res UserServer
	err := self.db.Where(&UserServer{UserID: userID, ServerID: serverID}).First(&res).Error
	return res, err
}

func (self *ChatApiDB) GetAllUsersByServerID(serverID uint) ([]User, error) {
	var server Server
	err := self.db.
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		First(&server, serverID).Error
	users := server.Users
	// for i, _ := range(users) {
	//     users[i].Password = ""
	// }
	return users, err
}

func (self *ChatApiDB) GetAllServersByUserID(userID uint) ([]Server, error) {
	var user User
	err := self.db.Preload("Servers").First(&user, userID).Error
	return user.Servers, err
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
