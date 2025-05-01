package repository

import (
	"context"
	"encoding/base64"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/database"
	"lembrago.com/lembrago/errors"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/utils"
)

func CreateVault(vault *models.Vault) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vaults")
	if vault.ID == primitive.NilObjectID {
		vault.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, vault)
	return err
}

func FindAllVaultsByOrgID(orgID primitive.ObjectID) ([]models.Vault, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vaults")
	cursor, err := collection.Find(ctx, bson.M{"orgId": orgID, "personalVault": false})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var vaults []models.Vault
	if err = cursor.All(ctx, &vaults); err != nil {
		return nil, err
	}

	return vaults, nil
}

func FindVaultByID(id primitive.ObjectID) (*models.Vault, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vaults")

	var vault models.Vault
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&vault)
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

func RemoveVaultByID(id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vaults")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func FindAllVaultsByUserOrgID(orgID primitive.ObjectID, userID primitive.ObjectID) ([]models.VaultWithMemberInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("vault_members")
	cursor, err := collection.Find(ctx, bson.M{"orgId": orgID, "userId": userID})
	if err != nil {
		return nil, err
	}

	var vaultMembers []models.VaultMember
	if err = cursor.All(ctx, &vaultMembers); err != nil {
		return nil, err
	}

	var vaults []models.VaultWithMemberInfo
	for _, vaultMember := range vaultMembers {
		vault, err := FindVaultByID(vaultMember.VaultID)
		vaultWithMemberInfo := models.VaultWithMemberInfo{
			ID:                     vault.ID.Hex(),
			OrgID:                  vault.OrgID.Hex(),
			EncryptedVaultMetadata: utils.FacEncryptedKeyDto(vault.EncryptedVaultMetadata.Ciphertext, vault.EncryptedVaultMetadata.Nonce),
			PersonalVault:          vault.PersonalVault,
			VaultCreatedBy:         vault.CreatedBy.Hex(),
			VaultUpdatedAt:         vault.UpdatedAt.Time().Format(time.RFC3339),
			VaultCreatedAt:         vault.CreatedAt.Time().Format(time.RFC3339),

			UserID:         vaultMember.UserID.Hex(),
			Permission:     string(vaultMember.Permission),
			ESVK_PubK_User: base64.StdEncoding.EncodeToString(vaultMember.ESVK_PubK_User),
			AddedBy:        vaultMember.AddedBy.Hex(),
			AddAt:          vaultMember.AddAt.Time().Format(time.RFC3339),
		}
		if err != nil {
			return nil, err
		}
		vaults = append(vaults, vaultWithMemberInfo)
	}

	return vaults, nil
}

func FindPasswordByID(id primitive.ObjectID) (*models.Password, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("passwords_items")

	var password models.Password
	err := collection.FindOne(ctx, bson.M{"_id": id}).Decode(&password)
	if err != nil {
		return nil, err
	}

	return &password, nil
}

func AddPasswordToVault(password models.Password) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("passwords_items")
	if password.ID == primitive.NilObjectID {
		password.ID = primitive.NewObjectID()
	}

	_, err := collection.InsertOne(ctx, password)
	return err
}

func RemovePasswordFromVault(passwordID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("passwords_items")
	_, err := collection.DeleteOne(ctx, bson.M{"_id": passwordID})
	return err
}

func UpdatePasswordInVault(
	passwordID primitive.ObjectID,
	newEncryptedData models.EncryptedKey,
	modifiedByID primitive.ObjectID,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("passwords_items")

	filter := bson.M{"_id": passwordID}

	updateFields := bson.M{
		"encryptedItemData": newEncryptedData,
		"lastModifiedBy":    modifiedByID,
		"updatedAt":         primitive.NewDateTimeFromTime(time.Now()),
	}
	updateDoc := bson.M{"$set": updateFields}

	result, err := collection.UpdateOne(ctx, filter, updateDoc)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.NewAppError(404, "Password not found")
	}

	return nil
}

func FindAllPasswordsByVaultID(vaultID primitive.ObjectID) ([]models.Password, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.GetCollection("passwords_items")
	cursor, err := collection.Find(ctx, bson.M{"vaultId": vaultID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var passwords []models.Password
	if err = cursor.All(ctx, &passwords); err != nil {
		return nil, err
	}

	return passwords, nil
}
