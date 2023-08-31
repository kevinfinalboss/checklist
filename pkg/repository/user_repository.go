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

func FindUserByCPF(hashedCPF string) (*models.User, error) {
	var user models.User
	collection := connection.Client.Database("checklist-apps").Collection("users")
	err := collection.FindOne(context.Background(), bson.M{"cpf": hashedCPF}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func FindAllUsers() ([]models.User, error) {
	var users []models.User
	collection := connection.Client.Database("checklist-apps").Collection("users")
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var user models.User
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
