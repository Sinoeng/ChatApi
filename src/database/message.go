package database


type Message struct {
    ID uint `gorm:"primaryKey"`
    UserID uint `gorm:"not null"`
    ServerID uint
    Payload string
}

func (self *ChatApiDB) NewMessage(payload string, userID uint, serverID uint) error {
    msg := Message{
        UserID: userID,
        ServerID: serverID,
        Payload: payload,
    }
    return self.db.Create(&msg).Error
}

func (self *ChatApiDB) DeleteMessage(messageID uint) error {
    return self.db.Delete(&Message{ID: messageID}).Error
}

func (self *ChatApiDB) GetMessageByID(messageID uint) (Message, error) {
    var res Message
    err := self.db.Where(&Message{ID: messageID}).Find(&res).Error
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
