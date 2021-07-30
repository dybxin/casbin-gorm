package model

import "time"

type User struct {
	ID        string     `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" sql:"index"`

	Name     string `json:"name" gorm:"unique; not null"`       //昵称
	Email    string `json:"email" gorm:"unique; not null"`      //邮箱地址
	Mobile   string `json:"mobile" gorm:"unique; not null"`     //手机号
	Password string `json:"password,omitempty" gorm:"not null"` //密码

	AuthorityID string    `json:"authority_id""`
	Authority   Authority `json:"authority"`
}

// TableName 表名
func (user *User) TableName() string {
	return "users"
}
