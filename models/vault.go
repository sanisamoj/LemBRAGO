package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Vault struct {
	ID                     primitive.ObjectID `bson:"_id"`
	OrgID                  primitive.ObjectID `bson:"orgId"`
	EncryptedVaultMetadata EncryptedKey       `bson:"encryptedVaultMetadata"`
	PersonalVault          bool               `bson:"personalVault"`
	CreatedBy              primitive.ObjectID `bson:"createdBy"`
	UpdatedAt              primitive.DateTime `bson:"updatedAt"`
	CreatedAt              primitive.DateTime `bson:"createdAt"`
}

type VaultMember struct {
	ID             primitive.ObjectID `bson:"_id"`
	VaultID        primitive.ObjectID `bson:"vaultId"`
	OrgID          primitive.ObjectID `bson:"orgId"`
	UserID         primitive.ObjectID `bson:"userId"`
	ESVK_PubK_User []byte             `bson:"esvk_pubK_user"`
	Permission     VaultPermission    `bson:"permission"`
	AddedBy        primitive.ObjectID `bson:"addedBy"`
	AddAt          primitive.DateTime `bson:"addAt"`
}

type VaultWithMemberInfo struct {
	ID                     string          `json:"id"`
	OrgID                  string          `json:"orgId"`                  // Vindo do Vault
	EncryptedVaultMetadata EncryptedKeyDto `json:"encryptedVaultMetadata"` // String base64 no json?
	PersonalVault          bool            `json:"personalVault"`
	VaultCreatedBy         string          `json:"vaultCreatedBy"`
	VaultUpdatedAt         string          `json:"vaultUpdatedAt"` // Ou time.Time
	VaultCreatedAt         string          `json:"vaultCreatedAt"` // Ou time.Time

	UserID         string `json:"userId"`
	Permission     string `json:"permission"`   // Ou o tipo que for
	ESVK_PubK_User string `json:"esvkPubKUser"` // Nome do campo como no json
	AddedBy        string `json:"addedBy"`
	AddAt          string `json:"addAt"` // Ou time.Time se preferir
}

type VaultPermission string

const (
	ADMIN VaultPermission = "admin"
	WRITE VaultPermission = "write"
	READ  VaultPermission = "read"
)
