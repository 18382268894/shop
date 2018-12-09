/**
*FileName: userModel
*Create on 2018/12/8 上午1:09
*Create by mok
*/

package model

const(
	//账号状态
	StatusUserDisActivated = 0
	StatusUserActivated = 1

)


type  User struct {
	UID uint64 `gorm:"primary_key;not null;"`
	Username string `gorm:"size:20;not null"`
	Password string `gorm:"size:50;not null"`
	Email string  `gorm:"not null;size:30;unique"`
	CreateTime int64 `gorm:"not null; default:0"`
	LoginIP int	`gorm:"not null;default:0"`
	LoginTime int64 `gorm:"not null;default:0"`
	Status uint8 `gorm:"default:0"`
}

func (u *User)TableName()string{
	return "shop_users"
}

