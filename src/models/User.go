package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TUser 用户表
const TUser = "t_user"

// User 用户schema
type User struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` // id
	Name     string             `json:"name" bson:"name"`                  // 用户昵称
	Pwd      string             `json:"pwd" bson:"pwd"`                    // 密码
	Email    string             `json:"email" bson:"email"`                //密码
	CreateAt time.Time          `json:"createAt" bson:"createAt"`          // 创建时间
}
