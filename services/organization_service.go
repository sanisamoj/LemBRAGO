package services

import (
	"encoding/base64"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lembrago.com/lembrago/cache"
	"lembrago.com/lembrago/errors"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/repository"
	"lembrago.com/lembrago/utils"
)

func CreateOrganization(request *models.CreateOrganizationRequest) (*models.Organization, error) {
	if err := utils.GetValidator().Struct(request); err != nil {
		return nil, err
	}

	id := primitive.NewObjectID()

	org := &models.Organization{
		ID:                 id,
		Name:               request.Name,
		ImageUrl:           request.ImageUrl,
		Email:              request.Email,
		SubscriptionPlan:   models.SubscriptionPlan(request.SubscriptionPlan),
		SubscriptionStatus: "active",
		UpdatedAt:          primitive.NewDateTimeFromTime(time.Now()),
		CreatedAt:          primitive.NewDateTimeFromTime(time.Now()),
	}

	saltPvBytes, err := base64.StdEncoding.DecodeString(request.User.PasswordVerifier.Salt)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 salt_pv format")
	}
	pvBytes, err := base64.StdEncoding.DecodeString(request.User.PasswordVerifier.Verifier)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 pv format")
	}

	saltEkBytes, err := base64.StdEncoding.DecodeString(request.User.Salt_ek)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 salt_ek format")
	}

	publicKey, err := base64.StdEncoding.DecodeString(request.User.Keys.PublicKey)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 pubKey format")
	}

	epk_cyphertext, err := base64.StdEncoding.DecodeString(request.User.Keys.EncryptedPrivateKey.Ciphertext)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 esk format")
	}
	epk_nonce, err := base64.StdEncoding.DecodeString(request.User.Keys.EncryptedPrivateKey.Nonce)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 salt_ek format")
	}

	esk_cyphertext, err := base64.StdEncoding.DecodeString(request.User.Keys.EncryptedSecretKey.Ciphertext)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 esk format")
	}
	esk_nonce, err := base64.StdEncoding.DecodeString(request.User.Keys.EncryptedSecretKey.Nonce)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 salt_ek format")
	}

	keys := &models.Keys{
		PublicKey: publicKey,
		EncryptedPrivateKey: models.EncryptedKey{
			Ciphertext: epk_cyphertext,
			Nonce:      epk_nonce,
		},
		EncryptedSecretKey: models.EncryptedKey{
			Ciphertext: esk_cyphertext,
			Nonce:      esk_nonce,
		},
	}

	user := &models.User{
		ID:               id,
		OrgID:            id,
		Username:         request.User.Username,
		Email:            request.Email,
		PasswordVerifier: pvBytes,
		Parameters:       request.User.PasswordVerifier.Parameters,
		SaltPV:           saltPvBytes,
		SaltEk:           saltEkBytes,
		Keys:             *keys,
		Role:             models.RoleAdmin,
		Status:           models.StatusActive,
		UpdatedAt:        primitive.NewDateTimeFromTime(time.Now()),
		CreatedAt:        primitive.NewDateTimeFromTime(time.Now()),
	}

	err = repository.CreateOrganization(org)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil, errors.NewAppError(409, err.Error())
		}
		return nil, err
	}
	repository.CreateUser(user)
	if request.User.MyVault != nil {
		_, err = CreatePersonalVault(id.Hex(), id.Hex(), request.User.MyVault)
		if err != nil {
			repository.DeleteOrganization(org.ID)
			repository.DeleteUser(id)
			return nil, err
		}
	}

	return org, nil
}

func InviteUser(adminID, orgID, email string, role models.UserRole) (string, error) {
	adminObjID, err := primitive.ObjectIDFromHex(adminID)
	if err != nil {
		return "", err
	}

	admin, err := repository.FindUserByID(adminObjID)
	if err != nil {
		return "", errors.NewAppError(403, "Unauthorized")
	}

	if admin.Role != models.RoleAdmin {
		return "", errors.NewAppError(403, "Only admins can invite users")
	}

	orgIbjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return "", errors.NewAppError(400, "Invalid OrgID")
	}

	user, err := repository.FindUserByEmailOrgID(email, orgIbjID)
	if user != nil {
		return "", errors.NewAppError(403, "User already exists")
	}

	organization, err := repository.FindOrganizationByID(orgIbjID)
	if err != nil {
		return "", errors.NewAppError(403, "Invalid Permission")
	}

	token, err := utils.GenerateJWTUserCreation(orgID, orgID, role, email)
	if err != nil {
		return "", err
	}
	code := primitive.NewObjectID().Hex()
	inviteCode := models.MinOrgWithTokenResponse{
		Token:        token,
		Organization: organization.Name,
		ImgUrl:       organization.ImageUrl,
		UserEmail:    email,
	}

	cache.SetStruct(code, &inviteCode)
	cache.SetTTL(code, 5*time.Minute)

	return code, nil
}

func GetInviteCodeToken(code string) (*models.MinOrgWithTokenResponse, error) {
	var minOrgWithTokenResponse models.MinOrgWithTokenResponse
	err := cache.GetStruct(code, &minOrgWithTokenResponse)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid code")
	}
	if (models.MinOrgWithTokenResponse{}) == minOrgWithTokenResponse {
		return nil, errors.NewAppError(400, "Invalid code")
	}

	return &minOrgWithTokenResponse, nil
}
