package database

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Email    string    `form:"default:null"`
	Username string    `gorm:"not null;unique"`
	Password string    `gorm:"not null"`
	Servers  []*Server `gorm:"many2many:user_server"`
}

func (self *ChatApiDB) InsertNewUser(username, password string) error {
	user := User{
		Username: username,
		Password: password,
	}
	return self.db.Create(&user).Error
}

func (self *ChatApiDB) InsertNewUserWEmail(username, password, email string) error {
	user := User{
		Username: username,
		Password: password,
		Email:    email,
	}
	return self.db.Create(&user).Error
}

func (self *ChatApiDB) ChangeEmailByID(userID uint, email string) error {
	tx := self.db.Begin()
	var user User
	err := tx.Where(&User{ID: userID}).First(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	user.Email = email
	err = tx.Save(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	return err
}

func (self *ChatApiDB) DeleteUserByID(id uint) error {
	return self.db.Delete(&User{ID: id}).Error
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

func (self *ChatApiDB) GetAllUsers() ([]User, error) {
	var res []User
	err := self.db.Find(&res).Error
	return res, err
}
