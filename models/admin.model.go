package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)
type Admin struct {
	ID            primitive.ObjectID `bson:"_id"`
	Name          *string            `json:"name" validate:"required" `
	Email         *string            `json:"email" validate:"required"`
	Contact       *string            `json:"contact" validate:"required"`
	Password      *string            `json:"password" validate:"required"`
	Company       *string            `json:"company" `
	Token         *string            `json:"token" `
	Refresh_token *string           ` json:"refresh_token"`
	User_id       *string            `json:"user_id" `
	User_Type     *string           ` json:"user_type"  validate:"required"`
}

type Request_to_admin struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title              *string            `json:"name" validate:"required"`
	Sendto             *string            `json:"sendto" validate:"required"`
	Discription        *string            `json:"discription" validate:"required"`
	Short_discription  *string            `json:"short_discription" validate:"required"`
	Sended_At          *string            `json:"sended_at" validate:"required"`
	Status_review      *string            `json:"status_review"`
	Status_reviewed_at *string            `json:"status_reviewed_at"`
}

type PDFUploads struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title     string             `json:"title" validate:"required"`
	UserEmail string             `json:"user_email" validate:"required"`
	PDFFile   []byte             `json:"pdf_file" validate:"required"`
}

type JobListing struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Role         string             `json:"role" bson:"role"`
	Location     string             `json:"location" bson:"location"`
	Company      string             `json:"company" bson:"company"`
	Description  string             `json:"desc" bson:"desc"`
	Requirements []string           `json:"requirements" bson:"requirements"`
	Link         string             `json:"link" bson:"link"`
	ResponsesLink string            `json:"responsesLink" bson:"responsesLink"`
	Active       bool               `json:"active" bson:"active"`
}
