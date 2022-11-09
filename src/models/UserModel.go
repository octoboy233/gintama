package models

import "fmt"

type UserModel struct {
	UserName string `gorm:"column:name"`                                //gorm
	UserId   int    `gorm:"column:id" uri:"id" binding:"required,gt=0"` //gin的用法
	Views    int    `gorm:"column:views"`
}

func NewUserModel() *UserModel {
	return &UserModel{}
}

func (this *UserModel) String() string {
	return fmt.Sprintf("user:id=%d,name=%s", this.UserId, this.UserName)
}
