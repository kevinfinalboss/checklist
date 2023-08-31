package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name              string             `bson:"name"`
	Email             string             `bson:"email"`
	Password          string             `bson:"password"`
	ConfirmPassword   string             `bson:"-"`
	CPF               string             `bson:"cpf"`
	BirthDate         string             `bson:"birthDate"`
	Address           string             `bson:"address"`
	HashedCPFForCheck string             `bson:"hashedCPFForCheck"`
}
