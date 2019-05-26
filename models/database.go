package models

import (
	"github.com/globalsign/mgo/bson"
)

type (
	User struct {
		ID        bson.ObjectId   `bson:"_id,omitempty" json:"id,omitempty"`
		Username  string          `bson:"username" json:"username"`
		Password  string          `bson:"password" json:"password,omitempty"`
		FirstName string          `bson:"firstname" json:"firstname"`
		LastName  string          `bson:"lastname" json:"lastname"`
		Email     string          `bson:"email" json:"email,omitempty"`
		Blogs     []bson.ObjectId `bson:"blogs,omitempty" json:"blogs,omitempty"`
		Token     string          `bson:"token,omitempty" json:"token,omitempty"`
		Picture   string          `bson:"profile_picture" json:"profile_picture,omitempty"`
	}
	Resp struct {
		Code   int         `json:"code"`
		Result interface{} `json:"result"`
	}
)
