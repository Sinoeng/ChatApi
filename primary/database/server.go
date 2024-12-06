package database

type Server struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"not null"`
	Users []*User `gorm:"many2many:user_server"`
}

func (self *ChatApiDB) AddUser(serverID uint, userID uint) error {
	var user User
	var server Server
	err := self.db.First(&user, userID).Error
	if err != nil {
		return err
	}
	err = self.db.First(&server, serverID).Error
	if err != nil {
		return err
	}
	err = self.db.Model(&server).Association("Users").Append(&user)
	return err
}
