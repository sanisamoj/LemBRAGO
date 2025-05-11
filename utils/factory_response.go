package utils

import (
	"time"

	"lembrago.com/lembrago/models"
)

func FacEncryptedKeyDto(ciphertext, nonce []byte) models.EncryptedKeyDto {
	return models.EncryptedKeyDto{
		Ciphertext: BytesToBase64(ciphertext),
		Nonce:      BytesToBase64(nonce),
	}
}

func FacVaultResponse(vault *models.Vault, email string, vaultMember *models.VaultMember) *models.VaultResponse {
	return &models.VaultResponse{
		ID:                     vault.ID.Hex(),
		OrgID:                  vault.OrgID.Hex(),
		EncryptedVaultMetadata: FacEncryptedKeyDto(vault.EncryptedVaultMetadata.Ciphertext, vault.EncryptedVaultMetadata.Nonce),
		MyMembership:           *FacVaultMemberResponse(email, vaultMember),
		PersonalVault:          vault.PersonalVault,
		CreatedBy:              vault.CreatedBy.Hex(),
		UpdatedAt:              vault.UpdatedAt.Time().Format(time.RFC3339),
		CreatedAt:              vault.CreatedAt.Time().Format(time.RFC3339),
	}
}

func FacVaultMemberResponse(email string, vaultMember *models.VaultMember) *models.VaultMemberResponse {
	return &models.VaultMemberResponse{
		ID:             vaultMember.ID.Hex(),
		VaultID:        vaultMember.VaultID.Hex(),
		UserID:         vaultMember.UserID.Hex(),
		Email:          email,
		Username:       nil,
		ESVK_PubK_User: BytesToBase64(vaultMember.ESVK_PubK_User),
		Permission:     string(vaultMember.Permission),
		AddedBy:        vaultMember.AddedBy.Hex(),
		AddAt:          vaultMember.AddAt.Time().Format(time.RFC3339),
	}
}

func FacMinimalUserRes(user *models.User) models.MinimalUserInfoResponse {
	return models.MinimalUserInfoResponse{
		ID:        user.ID.Hex(),
		OrgID:     user.OrgID.Hex(),
		Email:     user.Email,
		Username:  user.Username,
		PublicKey: BytesToBase64(user.Keys.PublicKey),
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Time().Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Time().Format(time.RFC3339),
	}
}
