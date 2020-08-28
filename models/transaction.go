package models

import (
	"github.com/jinzhu/gorm"
)

type (
	Transaction struct {
		gorm.Model
		State    string
		Amount    string
		TransactionId string
		Canceled bool
		User *User
	}
)
