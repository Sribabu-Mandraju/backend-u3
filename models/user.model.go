package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Firstname   string             `bson:"firstname" json:"firstname"`
	Lastname    string             `bson:"lastname" json:"lastname"`
	Email       string             `bson:"email" json:"email"`
	Password    string             `bson:"password" json:"password"`
	Company     string             `bson:"company" json:"company"`
	Address     string             `bson:"address" json:"address"`
	Country     string             `bson:"country" json:"country"`
	City        string             `bson:"city" json:"city"`
	Province    string             `bson:"province" json:"province"`
	Postal      string             `bson:"postal" json:"postal"`
	CreatedAt   time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type InputUser struct{
	Email       string 			   `json:"email" bson:"email"`
	Password    string 			   `json:"password" bson:"password"`
}

func (u *User) SetUserTimeStamps(){
	currentTime := time.Now();
	if(u.CreatedAt.IsZero()){
		u.CreatedAt = currentTime
	}
	u.UpdatedAt = currentTime
}