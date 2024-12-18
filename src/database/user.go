package database

type User struct {
	ID       uint     `gorm:"primaryKey"`
	Email    string   `form:"default:null"`
	Username string   `gorm:"not null;unique"`
	Password string   `gorm:"not null"`
	Servers  []Server `gorm:"many2many:user_servers;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:ServerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (self *ChatApiDB) InsertNewUser(username, password string) (uint, error) {
	user := User{
		Username: username,
		Password: password,
	}
	err := self.db.Create(&user).Error
	return user.ID, err
}

func (self *ChatApiDB) InsertNewUserWEmail(username, password, email string) (uint, error) {
	user := User{
		Username: username,
		Password: password,
		Email:    email,
	}
	err := self.db.Create(&user).Error
	return user.ID, err
}

func (self *ChatApiDB) ChangeUserEmailByID(userID uint, email string) error {
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
	if err != nil {
		tx.Rollback()
	}
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
