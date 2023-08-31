package repository

import (
	"context"

	connection "github.com/kevinfinalboss/checklist-apps/pkg/database"
	"github.com/kevinfinalboss/checklist-apps/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
)

func CreateUser(user *models.User) error {
	collection := connection.Client.Database("checklist-apps").Collection("users")
	_, err := collection.InsertOne(context.Background(), user)
	return err
}

func FindUserByEmail(email string) (*models.User, error) {
	var user models.User
	collection := connection.Client.Database("checklist-apps").Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindUserByHashedCPFForCheck(hashedCPFForCheck string) (*models.User, error) {
	var user models.User
	collection := connection.Client.Database("checklist-apps").Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"hashedCPFForCheck": hashedCPFForCheck}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
