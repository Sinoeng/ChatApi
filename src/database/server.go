package database

import "gorm.io/gorm/clause"

type Server struct {
	ID       uint      `gorm:"primaryKey"`
	Name     string    `gorm:"not null"`
	Users    []User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;many2many:user_servers;foreignKey:ID;joinForeignKey:ServerID;References:ID;joinReferences:UserID;"`
	Messages []Message `gorm:"foreignKey:ServerID"`
}

func (self *ChatApiDB) InsertNewServer(name string) (uint, error) {
	server := Server{
		Name: name,
	}
	err := self.db.Create(&server).Error
	return server.ID, err
}

func (self *ChatApiDB) GetServerByID(id uint) (Server, error) {
	var res Server
	err := self.db.Where(&Server{ID: id}).First(&res).Error
	return res, err
}

func (self *ChatApiDB) GetAllServers() ([]Server, error) {
	var res []Server
	err := self.db.Find(&res).Error
	return res, err
}

func (self *ChatApiDB) DeleteServerByID(serverID uint) error {
	return self.db.Select(clause.Associations).Delete(&Server{ID: serverID}).Error
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
