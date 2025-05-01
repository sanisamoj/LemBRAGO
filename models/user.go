package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID    primitive.ObjectID `bson:"_id" json:"id"`
	OrgID primitive.ObjectID `bson:"orgId" json:"orgId"`

	Username string `bson:"username" json:"username" validate:"required"`
	Email    string `bson:"email" json:"email" validate:"required,email"`

	PasswordVerifier []byte            `bson:"passwordVerifier" json:"passwordVerifier"` // PV
	SaltPV           []byte            `bson:"saltPV" json:"saltPV"`                     // salt_pv
	Parameters       Argo2IDParameters `bson:"parameters" json:"parameters"`

	SaltEk []byte `bson:"saltEk" json:"saltEk"` // salt_ek
	Keys   Keys   `bson:"keys" json:"keys"`

	Role   UserRole   `bson:"role", json:"role"`
	Status UserStatus `bson:"status"`

	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
	UpdatedAt primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
}

type Keys struct {
	PublicKey           []byte       `bson:"publicKey" json:"publicKey" validate:"required"`
	EncryptedPrivateKey EncryptedKey `bson:"encryptedPrivateKey" json:"encryptedPrivateKey" validate:"required"` // EUserPrivK
	EncryptedSecretKey  EncryptedKey `bson:"encryptedSecretKey" json:"encryptedSecretKey" validate:"required"`   // ESK
}

type EncryptedKey struct {
	Ciphertext []byte `bson:"ciphertext" json:"ciphertext" validate:"required"`
	Nonce      []byte `bson:"nonce" json:"nonce" validate:"required"`
}

type SavedMedia struct {
	ID       primitive.ObjectID `json:"id" bson:"_id"`
	OrgID    primitive.ObjectID `json:"orgId" bson:"orgId"`
	Filename string             `json:"filename" bson:"filename"`
	URL      string             `json:"url" bson:"url"`
	Size     int64              `json:"size" bson:"size"` // Bytes
	SavedAt  primitive.DateTime `json:"savedAt" bson:"savedAt"`
}

type UserRole string
type UserStatus string

const (
	RoleAdmin  UserRole = "admin"
	RoleMember UserRole = "member"

	StatusActive    UserStatus = "active"
	StatusInvited   UserStatus = "invited"
	StatusSuspended UserStatus = "suspended"
)
