package models

import (
	orm "gindriver/database"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	PublicKey string `json:"publickey"`
}

func (user User) Insert() (id int64, err error) {
	result := orm.Database.Create(&user)

	id = user.ID

	if result.Error != nil {
		err = result.Error
		return
	}

	return
}
