package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/database"
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

	return users, nil
}

func SaveMediaInRepo(svMedia *models.SavedMedia) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("saved_media")
	if svMedia.ID == primitive.NilObjectID {
		svMedia.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, svMedia)
	return err
}	

func GetMediaByID(id primitive.ObjectID) (*models.SavedMedia, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("saved_media")

	var svMedia models.SavedMedia
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&svMedia)
	if err != nil {
		return nil, err
	}

	return &svMedia, nil
}

func GetAllMediaByOrgID(orgID primitive.ObjectID) ([]models.SavedMedia, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("saved_media")
	cursor, err := collection.Find(ctx, bson.M{"orgId": orgID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var savedMedia []models.SavedMedia
	if err = cursor.All(ctx, &savedMedia); err != nil {
		return nil, err
	}

	return savedMedia, nil
}

func DeleteMediaInRepo(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("saved_media")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}