package domain

import (
	"gallery/server/utils"
	"github.com/jinzhu/gorm"
	"log"
)

type UserService interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(userId int) (*User, error)
}

func NewUserService(db *gorm.DB) *service {
	return &service{db}
}

type User struct {
	ID       int `gorm:"primary_key;auto_increment" json:"id"`
	Email    string `gorm:"size:100;not null;" json:"email"`
	Password string `gorm:"size:100;not null;" json:"password"`
}

func (r *service) SeedUsers() error {

	password, err := utils.Hash("password")
	if err != nil {
		return err
	}

	users := []User{
		{
			Email:    "steven@gmail.com",
			Password: string(password),
		},
		{
			Email:    "kenny@gmail.com",
			Password: string(password),
		},
	}

	for i := range users {
		err := r.db.Debug().Model(&User{}).Create(&users[i]).Error
		if err != nil {
			return err
		}
	}
	log.Printf("seedUsers routine OK !!!")
	return nil
}

func (r *service) GetUserByEmail(email string) (*User, error) {

	var user = &User{}
	err := r.db.Debug().Model(User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *service) GetUserByID(id int) (*User, error) {

	var user = &User{}
	err := r.db.Debug().Model(User{}).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}
