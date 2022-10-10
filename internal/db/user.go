package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Email       string    `json:"email" gorm:"email,unique"`
	Password    string    `json:"password" gorm:"password"`
	ConfirmedAt time.Time `json:"confirmedAt"`
}
