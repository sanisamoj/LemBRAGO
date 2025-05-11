package models

type VaultResponse struct {
	ID                     string              `json:"id"`
	OrgID                  string              `json:"orgId"`
	EncryptedVaultMetadata EncryptedKeyDto     `json:"encryptedVaultMetadata"`
	MyMembership           VaultMemberResponse `json:"myMembership"`
	PersonalVault          bool                `json:"personalVault"`
	CreatedBy              string              `json:"createdBy"`
	UpdatedAt              string              `json:"updatedAt"`
	CreatedAt              string              `json:"createdAt"`
}

type VaultMemberResponse struct {
	ID             string  `json:"id"`
	VaultID        string  `json:"vaultId"`
	UserID         string  `json:"userId"`
	Username       *string `json:"username,omitempty"`
	Email          string  `json:"email"`
	ESVK_PubK_User string  `json:"esvk_pubK_user"`
	Permission     string  `json:"permission"`
	AddedBy        string  `json:"addedBy"`
	AddAt          string  `json:"addAt"`
}

type CreateVaultRequest struct {
	EncryptedVaultMetadata EncryptedKeyDto `json:"e_vaultmetadata" validate:"required"`
	ESVK_PubK_User         string          `json:"esvk_pubK_user" validate:"required"`
	PersonalVault          *bool           `json:"personalVault"`
}

type UpdateVaultRequest struct {
	VaultId                string          `json:"vaultId" validate:"required"`
	EncryptedVaultMetadata EncryptedKeyDto `json:"e_vaultmetadata" validate:"required"`
	ESVK_PubK_User         string          `json:"esvk_pubK_user" validate:"required"`
}

type CreateVaultMemberRequest struct {
	VaultID        string          `json:"vaultId" validate:"required"`
	UserID         string          `json:"userId" validate:"required"`
	ESVK_PubK_User string          `json:"esvk_pubK_user" validate:"required"`
	Permission     VaultPermission `json:"permission" validate:"required,oneof=admin write read"`
}

type UpdateVaultMemberRequest struct {
	MemberID       string          `json:"memberId" validate:"required"`
	ESVK_PubK_User string          `json:"esvk_pubK_user" validate:"required"`
	Permission     VaultPermission `json:"permission" validate:"required,oneof=admin write read"`
}
