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

func AddVaultMember(vaultMember *models.VaultMember) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vault_members")
	if vaultMember.ID == primitive.NilObjectID {
		vaultMember.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, vaultMember)
	return err
}

func DeleteVaultMember(memberID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vault_members")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": memberID})
	return err
}

func UpdateVaultMember(memberID primitive.ObjectID, esvkBytes []byte, Permission models.VaultPermission) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	updateFields := bson.M{
		"esvk_pubK_user": esvkBytes,
		"permission":     Permission,
		"updatedAt":      time.Now(),
	}
	updateDoc := bson.M{"$set": updateFields}

	collection := database.GetCollection("vault_members")
	result, err := collection.UpdateOne(ctx, bson.M{"_id": memberID}, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.NewAppError(404, "Vault member not found")
	}

	return nil
}

func FindAllVaultMembersByVaultID(vaultID primitive.ObjectID) ([]models.VaultMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vault_members")
	cursor, err := collection.Find(ctx, bson.M{"vaultId": vaultID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vaultMembers []models.VaultMember
	if err = cursor.All(ctx, &vaultMembers); err != nil {
		return nil, err
	}

	return vaultMembers, nil
}

func FindMemberByUserVaultID(vaultId primitive.ObjectID, userID primitive.ObjectID) (*models.VaultMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vault_members")
	var vaultMember models.VaultMember
	err := collection.FindOne(ctx, bson.M{"vaultId": vaultId, "userId": userID}).Decode(&vaultMember)
	if err != nil {
		return nil, err
	}

	return &vaultMember, nil
}

func FindMemberById(vaultMemberID primitive.ObjectID) (*models.VaultMember, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vault_members")
	var vaultMember models.VaultMember
	err := collection.FindOne(ctx, bson.M{"_id": vaultMemberID}).Decode(&vaultMember)
	if err != nil {
		return nil, err
	}

	return &vaultMember, nil
}
