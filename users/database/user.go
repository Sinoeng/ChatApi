package database

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `form:"default:null"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

func (self *ChatApiDB) InsertNewUser(username, password string) error {
	user := User{
		Username: username,
		Password: password,
	}
	err := self.db.Create(&user).Error
	return err
}

func (self *ChatApiDB) GetUserByID(id uint) (User, error) {
    var res User
    err := self.db.Where(&User{ID: id}).First(&res).Error
    return res, err
}

func (self *ChatApiDB) GetUserByUsername(username string) (User, error) {
    var res User
    err := self.db.Where(&User{Username: username}).First(&res).Error
    return res, err
}

func (self *ChatApiDB) DeleteUserByID(id uint) error {
    return self.db.Delete(&User{ID: id}).Error
}

func (self *ChatApiDB) GetAllUsers() ([]User, error) {
	var res []User
	err := self.db.Find(&res).Error
	return res, err
}
