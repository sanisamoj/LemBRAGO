package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Password struct {
	ID                primitive.ObjectID `bson:"_id"`
	VaultID           primitive.ObjectID `bson:"vaultId"`
	EncryptedItemData EncryptedKey       `bson:"encryptedItemData"`
	CreatedBy         primitive.ObjectID `bson:"createdBy"`
	LastModifiedBy    primitive.ObjectID `bson:"lastModifiedBy"`
	CreatedAt         primitive.DateTime `bson:"createdAt"`
	UpdatedAt         primitive.DateTime `bson:"updatedAt"`
}

type CreatePasswordRequest struct {
	VaultID           string          `json:"vaultId" validate:"required"`
	EncryptedItemData EncryptedKeyDto `json:"encryptedItemData" validate:"required"`
}

type UpdatePasswordRequest struct {
	PasswordID        string          `json:"passwordId" validate:"required"`
	EncryptedItemData EncryptedKeyDto `json:"encryptedItemData" validate:"required"`
}

type PasswordResponse struct {
	ID                string          `json:"id"`
	VaultID           string          `json:"vaultId"`
	EncryptedItemData EncryptedKeyDto `json:"encryptedItemData"`
	CreatedAt         string          `json:"createdAt"`
	UpdatedAt         string          `json:"updatedAt"`
}
