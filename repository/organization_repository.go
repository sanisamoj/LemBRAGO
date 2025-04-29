package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/database"
	"lembrago.com/lembrago/models"
)

func CreateOrganization(org *models.Organization) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("organizations")
	if org.ID == primitive.NilObjectID {
		org.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, org)
	return err
}

func DeleteOrganization(orgID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("organizations")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": orgID})
	return err
}

func FindAllOrganizations() ([]models.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("organizations")
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var orgs []models.Organization
	if err = cursor.All(ctx, &orgs); err != nil {
		return nil, err
	}

	return orgs, nil
}

func FindOrganizationByID(id primitive.ObjectID) (*models.Organization, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("organizations")

	var org models.Organization
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&org)
	if err != nil {
		return nil, err
	}

	return &org, nil
}
