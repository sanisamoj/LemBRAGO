package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lembrago.com/lembrago/database"
	"lembrago.com/lembrago/models"
)

func RegisterAppVersion(version *models.ApplicationVersion) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("version")
	if version.ID == primitive.NilObjectID {
		version.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, version)
	return err
}

func GetAllVersions() ([]models.ApplicationVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("version")

	findOptions := options.Find()
	findOptions.SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var versions []models.ApplicationVersion
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func GetLastestVersion() (*models.ApplicationVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("version")
	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var version models.ApplicationVersion
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&version)
	if err != nil {
		return nil, err
	}

	return &version, nil
}

func UpdateAppVersion(version *models.ApplicationVersion) (*models.ApplicationVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("version")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": version.ID}, bson.M{"$set": version})
	if err != nil {
		return nil, err
	}

	return version, nil
}
