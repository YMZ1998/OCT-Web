package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username    string             `bson:"username" json:"username"`
	Password    string             `bson:"password,omitempty" json:"-"`
	Email       string             `bson:"email" json:"email"`
	Gender      string             `bson:"gender" json:"gender"`
	Age         int                `bson:"age" json:"age"`
	CreatedAt   int64              `bson:"created_at" json:"created_at"`
	LastLoginAt int64              `bson:"last_login_at,omitempty" json:"last_login_at,omitempty"`
}
