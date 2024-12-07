package database

type Server struct {
	ID    uint    `gorm:"primaryKey"`
	Name  string  `gorm:"not null"`
	Users []*User `gorm:"many2many:user_server"`
}

func (self *ChatApiDB) InsertNewServer(name string) error {
    server := Server{
        Name: name,
    }
    return self.db.Create(&server).Error
}

func (self *ChatApiDB) DeleteServer(serverID uint) error {
    return self.db.Delete(&Server{ID: serverID}).Error
}

func (self *ChatApiDB) ChangeServerName(serverID uint, newName string) error {
    tx := self.db.Begin()
    var server Server
    err := tx.Where(&Server{ID: serverID}).First(&server).Error
    if err != nil {
        tx.Rollback()
        return err
    }
    server.Name = newName
    err = tx.Save(&server).Error
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
