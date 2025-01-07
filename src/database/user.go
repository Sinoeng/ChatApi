package database

import "gorm.io/gorm/clause"

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Username string    `gorm:"not null;unique"`
	Password string    `gorm:"not null"`
	Email    UserEmail `gorm:"foreignKey:UserID"`
	Admin    bool      `gorm:"default:false"`
	Servers  []Server  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;many2many:user_servers;foreignKey:ID;joinForeignKey:UserID;References:ID;joinReferences:ServerID;"`
	Messages []Message `gorm:"foreignKey:UserID"`
}

type UserEmail struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint
	Email  string `gorm:"default:null"`
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
		Email: UserEmail{
			Email: email,
		},
	}
	err := self.db.Create(&user).Error
	return user.ID, err
}

func (self *ChatApiDB) InsertNewUserAsAdmin(username, password string) (uint, error) {
	user := User{
		Username: username,
		Password: password,
		Admin:    true,
	}
	err := self.db.Create(&user).Error
	return user.ID, err
}

func (self *ChatApiDB) InsertNewUserWEmailAsAdmin(username, password, email string) (uint, error) {
	user := User{
		Username: username,
		Password: password,
		Email:    UserEmail{
			Email: email,
		},
		Admin:    true,
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
	user.Email.Email = email
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

func (self *ChatApiDB) ChangeUserPasswordByID(userID uint, newPassword string) error {
    tx := self.db.Begin()
    var user User
    err := tx.First(&user, userID).Error
	if err != nil {
		tx.Rollback()
		return err
	}
    user.Password = newPassword
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
	return self.db.Select(clause.Associations).Delete(&User{ID: id}).Error
}

func (self *ChatApiDB) GetUserByID(id uint) (User, error) {
	var res User
	err := self.db.Preload("Email").Omit("password").Where(&User{ID: id}).First(&res).Error
	return res, err
}

func (self *ChatApiDB) GetUserByUsername(username string) (User, error) {
	var res User
	err := self.db.Preload("Email").Omit("password").Where(&User{Username: username}).First(&res).Error
	return res, err
}

func (self *ChatApiDB) GetUserPasswordByUsername(username string) (User, error) {
    var res User
    err := self.db.Where(&User{Username: username}).First(&res).Error
    return res, err
}

func (self *ChatApiDB) GetAllUsers() ([]User, error) {
	var res []User
	err := self.db.Preload("Email").Omit("password").Find(&res).Error
	return res, err
}

func (self *ChatApiDB) MakeUserAdmin(userID uint) error {
	var user User
	tx := self.db.Begin()
	err := tx.First(&user, userID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	user.Admin = true
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

func (self *ChatApiDB) UnMakeUserAdmin(userID uint) error {
	var user User
	tx := self.db.Begin()
	err := tx.First(&user, userID).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	user.Admin = false
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
