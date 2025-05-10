package services

import (
	"encoding/base64"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/errors"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/repository"
	"lembrago.com/lembrago/utils"
)

func CreateVault(userID string, req *models.CreateVaultRequest) (*models.VaultResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	user, err := repository.FindUserByID(userObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "User not found")
	}
	if user.Role != models.RoleAdmin {
		return nil, errors.NewAppError(403, "Only admin can create vault")
	}

	if req.PersonalVault != nil && *req.PersonalVault {
		return CreatePersonalVault(userID, user.OrgID.Hex(), req)
	}

	vaultID := primitive.NewObjectID()

	eskvBytes, err := utils.Base64ToBytes(req.ESVK_PubK_User)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid eskv")
	}

	evm_cyphertext, err := utils.Base64ToBytes(req.EncryptedVaultMetadata.Ciphertext)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid EncryptedVaultMetadata")
	}
	evm_nonce, err := utils.Base64ToBytes(req.EncryptedVaultMetadata.Nonce)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid EncryptedVaultMetadata")
	}

	encryptedVaultMetadata := models.EncryptedKey{
		Ciphertext: evm_cyphertext,
		Nonce:      evm_nonce,
	}

	vault := models.Vault{
		ID:                     vaultID,
		OrgID:                  user.OrgID,
		EncryptedVaultMetadata: encryptedVaultMetadata,
		PersonalVault:          false,
		CreatedBy:              user.ID,
		UpdatedAt:              primitive.NewDateTimeFromTime(time.Now()),
		CreatedAt:              primitive.NewDateTimeFromTime(time.Now()),
	}
	vaultMember := models.VaultMember{
		ID:             primitive.NewObjectID(),
		VaultID:        vaultID,
		OrgID:          user.OrgID,
		UserID:         user.ID,
		ESVK_PubK_User: eskvBytes,
		Permission:     models.ADMIN,
		AddedBy:        user.ID,
		AddAt:          primitive.NewDateTimeFromTime(time.Now()),
	}

	repository.CreateVault(&vault)
	repository.AddVaultMember(&vaultMember)

	vaultResponse := utils.FacVaultResponse(&vault, user.Email, &vaultMember)

	return vaultResponse, nil
}

func CreatePersonalVault(userID, orgID string, req *models.CreateVaultRequest) (*models.VaultResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	_, err = repository.FindVaultByID(userObjID)
	if err == nil {
		return nil, errors.NewAppError(403, "User already has a personal vault")
	}

	user, err := repository.FindUserByID(userObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "User not found")
	}

	orgObjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid orgID")
	}

	eskvBytes, err := utils.Base64ToBytes(req.ESVK_PubK_User)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid eskv")
	}

	evm_cyphertext, err := utils.Base64ToBytes(req.EncryptedVaultMetadata.Ciphertext)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid EncryptedVaultMetadata")
	}
	evm_nonce, err := utils.Base64ToBytes(req.EncryptedVaultMetadata.Nonce)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid EncryptedVaultMetadata")
	}

	encryptedVaultMetadata := models.EncryptedKey{
		Ciphertext: evm_cyphertext,
		Nonce:      evm_nonce,
	}

	vault := models.Vault{
		ID:                     userObjID,
		OrgID:                  orgObjID,
		EncryptedVaultMetadata: encryptedVaultMetadata,
		PersonalVault:          true,
		CreatedBy:              userObjID,
		UpdatedAt:              primitive.NewDateTimeFromTime(time.Now()),
		CreatedAt:              primitive.NewDateTimeFromTime(time.Now()),
	}
	vaultMember := models.VaultMember{
		ID:             primitive.NewObjectID(),
		VaultID:        userObjID,
		OrgID:          orgObjID,
		UserID:         userObjID,
		ESVK_PubK_User: eskvBytes,
		Permission:     models.ADMIN,
		AddedBy:        userObjID,
		AddAt:          primitive.NewDateTimeFromTime(time.Now()),
	}

	repository.CreateVault(&vault)
	repository.AddVaultMember(&vaultMember)

	encryptedVaultMetadataRes := utils.FacEncryptedKeyDto(evm_cyphertext, evm_nonce)
	vaultResponse := models.VaultResponse{
		ID:                     vault.ID.Hex(),
		OrgID:                  vault.OrgID.Hex(),
		EncryptedVaultMetadata: encryptedVaultMetadataRes,
		MyMembership: models.VaultMemberResponse{
			ID:             vaultMember.ID.Hex(),
			VaultID:        vaultMember.VaultID.Hex(),
			UserID:         vaultMember.UserID.Hex(),
			Email:          user.Email,
			ESVK_PubK_User: utils.BytesToBase64(vaultMember.ESVK_PubK_User),
			Permission:     string(vaultMember.Permission),
			AddedBy:        vaultMember.AddedBy.Hex(),
			AddAt:          vaultMember.AddAt.Time().Format(time.RFC3339),
		},
		PersonalVault: true,
		CreatedBy:     vault.CreatedBy.Hex(),
		UpdatedAt:     vault.UpdatedAt.Time().Format(time.RFC3339),
		CreatedAt:     vault.CreatedAt.Time().Format(time.RFC3339),
	}

	return &vaultResponse, nil
}

func RemoveVault(userID, vaultID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewAppError(400, "Invalid userID")
	}

	vaultObjID, err := primitive.ObjectIDFromHex(vaultID)
	if err != nil {
		return errors.NewAppError(400, "Invalid vaultID")
	}

	vault, err := repository.FindVaultByID(vaultObjID)
	if err != nil {
		return errors.NewAppError(404, "Vault not found")
	}

	vaultMem, err := repository.FindMemberByUserVaultID(vaultObjID, userObjID)
	if err != nil {
		return errors.NewAppError(403, "Vault member not found")
	}

	if vault.CreatedBy != userObjID || vaultMem.Permission != models.ADMIN {
		return errors.NewAppError(403, "You are not allowed to remove this vault")
	}

	err = repository.RemoveVaultByID(vaultObjID)
	if err != nil {
		return errors.NewAppError(500, "Failed to remove vault")
	}

	go removeVaultData(vaultObjID)
	return nil
}

func removeVaultData(vaultID primitive.ObjectID) error {
	members, err := repository.FindAllVaultMembersByVaultID(vaultID)
	if err != nil {
		return fmt.Errorf("Failed to find vault members: %v", err)
	}

	for _, member := range members {
		err = repository.DeleteVaultMember(member.ID)
		if err != nil {
			return fmt.Errorf("Failed to remove vault member: %v", err)
		}
	}

	passwords, err := repository.FindAllPasswordsByVaultID(vaultID)
	if err != nil {
		return fmt.Errorf("Failed to find passwords: %v", err)
	}

	for _, password := range passwords {
		err = repository.RemovePasswordFromVault(password.ID)
		if err != nil {
			return fmt.Errorf("Failed to remove password: %v", err)
		}
	}

	return nil
}

func FindMyAllVaultsByOrgID(userID, orgID string) ([]models.VaultWithMemberInfo, error) {
	orgObjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid orgID")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	vaultsWithMember, err := repository.FindAllVaultsByUserOrgID(orgObjID, userObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "Vaults not found")
	}

	if len(vaultsWithMember) == 0 {
		return nil, errors.NewAppError(404, "Vaults not found")
	}

	return vaultsWithMember, nil
}

func AddMemberToVault(userID string, req *models.CreateVaultMemberRequest) (*models.VaultMemberResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	vaultObjID, err := primitive.ObjectIDFromHex(req.VaultID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid vaultID")
	}

	permission, err := repository.FindMemberByUserVaultID(vaultObjID, userObjID)
	if err != nil {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}
	if permission.Permission != models.ADMIN {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}

	eskvBytes, err := base64.StdEncoding.DecodeString(req.ESVK_PubK_User)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid eskv")
	}

	memberObJID, err := primitive.ObjectIDFromHex(req.UserID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}
	vaultMember := models.VaultMember{
		ID:             primitive.NewObjectID(),
		VaultID:        vaultObjID,
		OrgID:          permission.OrgID,
		UserID:         memberObJID,
		ESVK_PubK_User: eskvBytes,
		Permission:     req.Permission,
		AddedBy:        userObjID,
		AddAt:          primitive.NewDateTimeFromTime(time.Now()),
	}

	repository.AddVaultMember(&vaultMember)

	targetUser, err := repository.FindUserByID(memberObJID)
	if err != nil {
		return nil, errors.NewAppError(404, "User not found")
	}

	vaultMemberResponse := models.VaultMemberResponse{
		ID:             vaultMember.ID.Hex(),
		VaultID:        vaultMember.VaultID.Hex(),
		UserID:         vaultMember.UserID.Hex(),
		Username:       targetUser.Username,
		Email:          targetUser.Email,
		ESVK_PubK_User: base64.StdEncoding.EncodeToString(vaultMember.ESVK_PubK_User),
		Permission:     string(vaultMember.Permission),
		AddedBy:        vaultMember.AddedBy.Hex(),
		AddAt:          vaultMember.AddAt.Time().Format(time.RFC3339),
	}

	return &vaultMemberResponse, nil
}

func RemoveMemberFromVault(userID, memberId string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewAppError(400, "Invalid userID")
	}

	memberObjId, err := primitive.ObjectIDFromHex(memberId)
	if err != nil {
		return errors.NewAppError(400, "Invalid vaultID")
	}

	vaultMember, err := repository.FindMemberById(memberObjId)
	if err != nil {
		return errors.NewAppError(404, "Member not found")
	}

	currentVaultMember, err := repository.FindMemberByUserVaultID(vaultMember.VaultID, userObjID)
	if err != nil {
		return errors.NewAppError(403, "Invalid Permission")
	}
	if currentVaultMember.Permission != models.ADMIN {
		return errors.NewAppError(403, "Invalid Permission")
	}

	err = repository.DeleteVaultMember(memberObjId)
	if err != nil {
		return errors.NewAppError(404, "Member not found")
	}

	return nil
}

func UpdateMemberPermission(userID, OrgID string, req *models.UpdateVaultMemberRequest) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewAppError(400, "Invalid userID")
	}

	orgObjID, err := primitive.ObjectIDFromHex(OrgID)
	if err != nil {
		return errors.NewAppError(400, "Invalid orgID")
	}

	permission, err := repository.FindMemberByUserVaultID(orgObjID, userObjID)
	if err != nil {
		return errors.NewAppError(403, "Invalid Permission")
	}
	if permission.Permission != models.ADMIN {
		return errors.NewAppError(403, "Invalid Permission")
	}

	memberObjId, err := primitive.ObjectIDFromHex(req.MemberID)
	if err != nil {
		return errors.NewAppError(400, "Invalid vaultID")
	}

	_, err = repository.FindMemberById(memberObjId)
	if err != nil {
		return errors.NewAppError(404, "Member not found")
	}

	esvkBytes, err := base64.StdEncoding.DecodeString(req.ESVK_PubK_User)
	if err != nil {
		return errors.NewAppError(400, "Invalid eskv")
	}

	err = repository.UpdateVaultMember(memberObjId, esvkBytes, req.Permission)
	if err != nil {
		return errors.NewAppError(404, "Member not found")
	}

	return nil
}

func GetAllMembersFromTheVault(vaultID, userID string) ([]models.VaultMemberResponse, error) {
	vaultObjID, err := primitive.ObjectIDFromHex(vaultID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid vaultID")
	}

	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	user, err := repository.FindUserByID(userObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "User not found")
	}
	if user.Role != models.RoleAdmin {
		return nil, errors.NewAppError(403, "Only admin can create vault")
	}

	vaultMembers, err := repository.FindAllVaultMembersByVaultID(vaultObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "Members not found")
	}

	if len(vaultMembers) == 0 {
		return nil, errors.NewAppError(404, "Members not found")
	}

	var vaultMemberResponses []models.VaultMemberResponse
	for _, vaultMember := range vaultMembers {
		user, err := repository.FindUserByID(vaultMember.UserID)
		if err != nil {
			return nil, errors.NewAppError(500, "Unknown Error")
		}
		vaultMemberResponse := models.VaultMemberResponse{
			ID:             vaultMember.ID.Hex(),
			VaultID:        vaultMember.VaultID.Hex(),
			UserID:         vaultMember.UserID.Hex(),
			Username:       user.Username,
			Email:          user.Email,
			ESVK_PubK_User: base64.StdEncoding.EncodeToString(vaultMember.ESVK_PubK_User),
			Permission:     string(vaultMember.Permission),
			AddedBy:        vaultMember.AddedBy.Hex(),
			AddAt:          vaultMember.AddAt.Time().Format(time.RFC3339),
		}
		vaultMemberResponses = append(vaultMemberResponses, vaultMemberResponse)
	}

	return vaultMemberResponses, nil
}

func AddPasswordToVault(userID string, req *models.CreatePasswordRequest) (*models.PasswordResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	vaultObjID, err := primitive.ObjectIDFromHex(req.VaultID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid vaultID")
	}
	permission, err := repository.FindMemberByUserVaultID(vaultObjID, userObjID)
	if err != nil {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}
	if permission.Permission != models.ADMIN {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}

	cipherBytes, err := utils.Base64ToBytes(req.EncryptedItemData.Ciphertext)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 Ciphertext format")
	}
	nonceBytes, err := utils.Base64ToBytes(req.EncryptedItemData.Nonce)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 Nonce format")
	}

	eid := models.EncryptedKey{
		Ciphertext: cipherBytes,
		Nonce:      nonceBytes,
	}
	password := models.Password{
		ID:                primitive.NewObjectID(),
		VaultID:           vaultObjID,
		EncryptedItemData: eid,
		CreatedBy:         userObjID,
		LastModifiedBy:    userObjID,
		CreatedAt:         primitive.NewDateTimeFromTime(time.Now()),
		UpdatedAt:         primitive.NewDateTimeFromTime(time.Now()),
	}

	err = repository.AddPasswordToVault(password)
	if err != nil {
		return nil, errors.NewAppError(500, "Unknown Error")
	}

	passwords := models.PasswordResponse{
		ID:      password.ID.Hex(),
		VaultID: password.VaultID.Hex(),
		EncryptedItemData: models.EncryptedKeyDto{
			Ciphertext: utils.BytesToBase64(password.EncryptedItemData.Ciphertext),
			Nonce:      utils.BytesToBase64(password.EncryptedItemData.Nonce),
		},
		CreatedAt: password.CreatedAt.Time().Format(time.RFC3339),
		UpdatedAt: password.UpdatedAt.Time().Format(time.RFC3339),
	}

	return &passwords, nil
}

func DeletePasswordFromVault(userID, passwordID string) error {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.NewAppError(400, "Invalid userID")
	}

	passwordObjID, err := primitive.ObjectIDFromHex(passwordID)
	if err != nil {
		return errors.NewAppError(400, "Invalid passwordID")
	}

	password, err := repository.FindPasswordByID(passwordObjID)
	if err != nil {
		return errors.NewAppError(404, "Password not found")
	}

	permission, err := repository.FindMemberByUserVaultID(password.VaultID, userObjID)
	if err != nil {
		return errors.NewAppError(403, "Invalid Permission")
	}
	if permission.Permission != models.WRITE && permission.Permission != models.ADMIN {
		return errors.NewAppError(403, "Invalid Permission")
	}

	err = repository.RemovePasswordFromVault(passwordObjID)
	if err != nil {
		return errors.NewAppError(500, "Unknown Error")
	}

	return nil
}

func UpdatePasswordInVault(userID string, req *models.UpdatePasswordRequest) (*models.PasswordResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	passwordObjID, err := primitive.ObjectIDFromHex(req.PasswordID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid passwordID")
	}

	password, err := repository.FindPasswordByID(passwordObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "Password not found")
	}

	permission, err := repository.FindMemberByUserVaultID(password.VaultID, userObjID)
	if err != nil {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}
	if permission.Permission != models.WRITE && permission.Permission != models.ADMIN {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}

	cipherBytes, err := utils.Base64ToBytes(req.EncryptedItemData.Ciphertext)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 Ciphertext format")
	}
	nonceBytes, err := utils.Base64ToBytes(req.EncryptedItemData.Nonce)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 Nonce format")
	}

	eid := models.EncryptedKey{
		Ciphertext: cipherBytes,
		Nonce:      nonceBytes,
	}
	password.EncryptedItemData = eid

	err = repository.UpdatePasswordInVault(passwordObjID, password.EncryptedItemData, userObjID)
	if err != nil {
		return nil, errors.NewAppError(500, "Unknown Error")
	}

	pRes := NewPasswordResponse(*password)

	return &pRes, nil
}

func GetAllPasswordsFromVault(userID, vaultID string) ([]models.PasswordResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID")
	}

	vaultObjID, err := primitive.ObjectIDFromHex(vaultID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid vaultID")
	}

	permission, err := repository.FindMemberByUserVaultID(vaultObjID, userObjID)
	if err != nil {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}
	if permission.Permission != models.READ && permission.Permission != models.WRITE && permission.Permission != models.ADMIN {
		return nil, errors.NewAppError(403, "Invalid Permission")
	}

	allEid, err := repository.FindAllPasswordsByVaultID(vaultObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "Passwords not found")
	}

	if len(allEid) == 0 {
		return nil, errors.NewAppError(404, "Passwords not found")
	}

	var passwords []models.PasswordResponse
	for _, password := range allEid {
		passwords = append(passwords, models.PasswordResponse{
			ID:      password.ID.Hex(),
			VaultID: password.VaultID.Hex(),
			EncryptedItemData: models.EncryptedKeyDto{
				Ciphertext: utils.BytesToBase64(password.EncryptedItemData.Ciphertext),
				Nonce:      utils.BytesToBase64(password.EncryptedItemData.Nonce),
			},
			CreatedAt: password.CreatedAt.Time().Format(time.RFC3339),
			UpdatedAt: password.UpdatedAt.Time().Format(time.RFC3339),
		})
	}

	return passwords, nil
}

func NewPasswordResponse(password models.Password) models.PasswordResponse {
	return models.PasswordResponse{
		ID:      password.ID.Hex(),
		VaultID: password.VaultID.Hex(),
		EncryptedItemData: models.EncryptedKeyDto{
			Ciphertext: utils.BytesToBase64(password.EncryptedItemData.Ciphertext),
			Nonce:      utils.BytesToBase64(password.EncryptedItemData.Nonce),
		},
		CreatedAt: password.CreatedAt.Time().Format(time.RFC3339),
		UpdatedAt: password.UpdatedAt.Time().Format(time.RFC3339),
	}
}
