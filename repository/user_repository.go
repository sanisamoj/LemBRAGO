package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/database"
	"lembrago.com/lembrago/errors"
	"lembrago.com/lembrago/models"
)

func CreateUser(user *models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("users")
	if user.ID == primitive.NilObjectID {
		user.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, user)

	return err
}

func DeleteUser(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("users")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func FindUserByEmailOrgID(email string, OrgID primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"email": email, "orgId": OrgID}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindAllUsersByEmail(email string) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("users")
	cursor, err := collection.Find(ctx, bson.M{"email": email})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func FindUserByID(id primitive.ObjectID) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUsersByOrgID(orgID primitive.ObjectID) ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("users")
	cursor, err := collection.Find(ctx, bson.M{"orgId": orgID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.NewAppError(404, "No users found")
	}

	return users, nil
}
