
package model //指明此包为 model 包

import (
	//The fantastic ORM library for Golang
	//一个非常好用的 go 语言 ORM 库
	"github.com/jinzhu/gorm"
)

//定义一个 User 的结构体, 用来存放用户名和密码
type User struct {
	gorm.Model        //加入此行用于在数据库中创建记录的 mate 数据
	UserID     string `gorm:"type:varchar(30);UNIQUE;unique_index;not null"`
	Password   string `gorm:"size:255"`
}

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}