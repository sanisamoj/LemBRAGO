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

func RegisterAppVersion(version *models.AppVersion) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("appVersion")
	if version.ID == primitive.NilObjectID {
		version.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, version)
	return err
}

func GetAllAppVersion() ([]models.AppVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("appVersion")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var versions []models.AppVersion
	if err = cursor.All(ctx, &versions); err != nil {
		return nil, err
	}

	return versions, nil
}

func GetLastestAppVersion() (*models.AppVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("appVersion")
	opts := options.FindOne().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	var version models.AppVersion
	err := collection.FindOne(ctx, bson.M{}, opts).Decode(&version)
	if err != nil {
		return nil, err
	}
	
	return &version, nil
}

func UpdateAppVersion(version *models.AppVersion) (*models.AppVersion, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("appVersion")
	_, err := collection.UpdateOne(ctx, bson.M{"_id": version.ID}, bson.M{"$set": version})
	if err != nil {
		return nil, err
	}

	return version, nil
}
