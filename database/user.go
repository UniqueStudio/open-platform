package database

import "gorm.io/gorm"

// FROM SSO database
// it's highly recommended not to modify it
type User struct {
	UID  string `json:"uid" gorm:"column:uid;primaryKey"`
	Role int    `json:"role" gorm:"column:role"`
}

func (u *User) TableName() string {
	return "user"
}

func GetUserByUID(uid string) (*User, error) {
	user := new(User)
	result := UserDB.Table(user.TableName()).Where("uid = ?", uid).Scan(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected < 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return user, nil
}
