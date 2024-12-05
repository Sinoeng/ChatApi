package database

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `form:"default:null"`
	Username string `gorm:"not null;unique"`
	Password string `gorm:"not null"`
}

func (self *UserDB) InsertNewUser(username, password string) error {
	user := User{
		Username: username,
		Password: password,
	}
	return self.db.Create(&user).Error
}

func (self *UserDB) InsertNewUserWEmail(username, password, email string) error {
    user := User{
        Username: username,
        Password: password,
        Email: email,
    }
    return self.db.Create(&user).Error
}

func (self *UserDB) DeleteUserByID(id uint) error {
    return self.db.Delete(&User{ID: id}).Error
}

func (self *UserDB) GetUserByID(id uint) (User, error) {
	var res User
	err := self.db.Where(&User{ID: id}).First(&res).Error
	return res, err
}

func (self *UserDB) GetUserByUsername(username string) (User, error) {
	var res User
	err := self.db.Where(&User{Username: username}).First(&res).Error
	return res, err
}

func (self *UserDB) GetAllUsers() ([]User, error) {
	var res []User
	err := self.db.Find(&res).Error
	return res, err
}
