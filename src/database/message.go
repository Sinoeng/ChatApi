package database

type Message struct {
	ID       uint `gorm:"primaryKey"`
	UserID   uint `gorm:"not null"`
	ServerID uint `gorm:"not null"`
	Payload  string
}

func (self *ChatApiDB) NewMessage(payload string, userID uint, serverID uint) (uint, error) {
	msg := Message{
		UserID:   userID,
		ServerID: serverID,
		Payload:  payload,
	}
    err := self.db.Create(&msg).Error
    return msg.ID, err
}

func (self *ChatApiDB) DeleteMessage(messageID uint) error {
	return self.db.Delete(&Message{ID: messageID}).Error
}

func (self *ChatApiDB) GetMessageByID(messageID uint) (Message, error) {
	var res Message
	err := self.db.Where(&Message{ID: messageID}).First(&res).Error
	return res, err
}

func (self *ChatApiDB) GetMessagesByServerID(serverID uint) ([]Message, error) {
	var res []Message
	err := self.db.Where(&Message{ServerID: serverID}).Find(&res).Error
	return res, err
}

func (self *ChatApiDB) GetMessagesByUserID(userID uint) ([]Message, error) {
	var res []Message
	err := self.db.Where(&Message{UserID: userID}).Find(&res).Error
	return res, err
}
