package services

import (
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"lembrago.com/lembrago/cache"
	"lembrago.com/lembrago/errors"
	"lembrago.com/lembrago/internal/config"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/repository"
	"lembrago.com/lembrago/utils"
)

const (
	uploadDir = "./uploads"
)

func SendAuthCode(email string) error {
	users, err := repository.FindAllUsersByEmail(email)
	if err != nil || len(users) == 0 {
		return errors.NewAppError(401, "Invalid Credentials")
	}

	code := utils.Gen6DigCod()
	go utils.SendAuthCodeEmail(email, code)
	go regAuthCode(email, code)
	return nil
}

func regAuthCode(email, code string) {
	key := fmt.Sprintf("auth-%s", email)
	cache.Set(key, code)
	cache.SetTTL(key, 5*time.Minute)
}

func GetLoginInfoFromUser(email, code string) ([]models.UserWithOrganizationResponse, error) {
	users, err := repository.FindAllUsersByEmail(email)
	if err != nil {
		return nil, errors.NewAppError(401, "User not found")
	}

	if len(users) == 0 {
		return nil, errors.NewAppError(401, "User not found")
	}

	key := fmt.Sprintf("auth-%s", email)
	cacheCode, err := cache.Get(key)
	if err != nil || code != cacheCode {
		return nil, errors.NewAppError(403, "Invalid Code")
	}

	cache.Delete(key)

	attKey := fmt.Sprintf("att-%s", email)
	cache.Delete(attKey)

	var userWithOrganizationResponseList []models.UserWithOrganizationResponse
	for _, user := range users {
		organization, err := repository.FindOrganizationByID(user.OrgID)
		if err != nil {
			return nil, errors.NewAppError(401, "User not found")
		}
		userWithOrganizationResponse := models.UserWithOrganizationResponse{
			OrgID:            user.OrgID.Hex(),
			OrganizationName: organization.Name,
			OrgImagUrl:       organization.ImageUrl,
			PasswordVerifier: models.PasswordVerifierResponse{
				Salt:       base64.StdEncoding.EncodeToString(user.SaltPV),
				Parameters: user.Parameters,
			},
		}
		userWithOrganizationResponseList = append(userWithOrganizationResponseList, userWithOrganizationResponse)
	}

	return userWithOrganizationResponseList, nil
}

func UserLogin(comparison *models.UserLoginComparison) (*models.UserLoginResponse, error) {
	orgObjID, err := primitive.ObjectIDFromHex(comparison.OrgID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid orgID")
	}

	user, err := repository.FindUserByEmailOrgID(comparison.Email, orgObjID)
	if err != nil {
		return nil, errors.NewAppError(401, "User not found")
	}

	verifierBytes, err := base64.StdEncoding.DecodeString(comparison.Verifier)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid base64 verifier format")
	}

	if subtle.ConstantTimeCompare(user.PasswordVerifier, verifierBytes) != 1 {
		return nil, errors.NewAppError(401, "Invalid credentials")
	}

	tokenStr, err := utils.GenerateJWT(user.ID.Hex(), user.OrgID.Hex(), user.Role)
	if err != nil {
		return nil, errors.NewAppError(500, "Unknonw Error")
	}

	userRespose := models.UserResponse{
		ID:       user.ID.Hex(),
		Email:    user.Email,
		OrgId:    user.OrgID.Hex(),
		Username: user.Username,

		PasswordVerifier: models.PasswordVerifierResponse{
			Salt:       base64.StdEncoding.EncodeToString(user.SaltPV),
			Parameters: user.Parameters,
		},

		Salt_ek: base64.StdEncoding.EncodeToString(user.SaltEk),

		Keys: models.KeysDTO{
			PublicKey:           utils.BytesToBase64(user.Keys.PublicKey),
			EncryptedPrivateKey: utils.FacEncryptedKeyDto(user.Keys.EncryptedPrivateKey.Ciphertext, user.Keys.EncryptedPrivateKey.Nonce),
			EncryptedSecretKey:  utils.FacEncryptedKeyDto(user.Keys.EncryptedSecretKey.Ciphertext, user.Keys.EncryptedSecretKey.Nonce),
		},
	}

	return &models.UserLoginResponse{User: userRespose, Token: tokenStr}, nil
}

func UserRegister(request *models.CreateUserRequest, orgID, email string, role models.UserRole) error {
	OrgObjectID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return errors.NewAppError(400, "Invalid orgID format")
	}

	saltPvBytes, err := utils.Base64ToBytes(request.PasswordVerifier.Salt)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 salt_pv format")
	}
	pvBytes, err := utils.Base64ToBytes(request.PasswordVerifier.Verifier)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 pv format")
	}
	saltEkBytes, err := utils.Base64ToBytes(request.Salt_ek)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 salt_ek format")
	}

	publicKey, err := utils.Base64ToBytes(request.Keys.PublicKey)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 pubKey format")
	}

	epk_cyphertext, err := utils.Base64ToBytes(request.Keys.EncryptedPrivateKey.Ciphertext)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 esk format")
	}
	epk_nonce, err := utils.Base64ToBytes(request.Keys.EncryptedPrivateKey.Nonce)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 salt_ek format")
	}

	esk_cyphertext, err := utils.Base64ToBytes(request.Keys.EncryptedSecretKey.Ciphertext)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 esk format")
	}
	esk_nonce, err := utils.Base64ToBytes(request.Keys.EncryptedSecretKey.Nonce)
	if err != nil {
		return errors.NewAppError(400, "Invalid base64 salt_ek format")
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

	userID := primitive.NewObjectID()
	user := &models.User{
		ID:               userID,
		OrgID:            OrgObjectID,
		Username:         request.Username,
		Email:            email,
		PasswordVerifier: pvBytes,
		Parameters:       request.PasswordVerifier.Parameters,
		SaltPV:           saltPvBytes,
		SaltEk:           saltEkBytes,
		Keys:             *keys,
		Role:             role,
		Status:           models.StatusActive,
		UpdatedAt:        primitive.NewDateTimeFromTime(time.Now()),
		CreatedAt:        primitive.NewDateTimeFromTime(time.Now()),
	}

	err = repository.CreateUser(user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.NewAppError(409, err.Error())
		}
		return err
	}

	if request.MyVault != nil {
		_, err = CreatePersonalVault(userID.Hex(), orgID, request.MyVault)
		if err != nil {
			repository.DeleteUser(userID)
			return err
		}
	}

	go utils.SendWelcomeEmail(user.Email, user.Username)

	return nil
}

func GetUserByID(userID string) (*models.MinimalUserInfoResponse, error) {
	userObjID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid userID format")
	}

	user, err := repository.FindUserByID(userObjID)
	if err != nil {
		return nil, errors.NewAppError(404, "User not found")
	}

	minimalUser := utils.FacMinimalUserRes(user)

	return &minimalUser, nil
}

func GetUsersByOrgID(orgID string) ([]models.MinimalUserInfoResponse, error) {
	OrgObjectID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, errors.NewAppError(400, "Invalid orgID format")
	}

	users, err := repository.FindUsersByOrgID(OrgObjectID)
	if err != nil {
		return nil, errors.NewAppError(404, "Users not found")
	}

	var minimalUsers []models.MinimalUserInfoResponse

	for _, user := range users {
		minimalUser := utils.FacMinimalUserRes(&user)
		minimalUsers = append(minimalUsers, minimalUser)
	}

	return minimalUsers, nil
}

func SignOut(token string) error {
	claims, err := utils.GetTokenInfo(config.GetServerConfig().JWTSecret, token)
	if err != nil {
		return err
	}

	err = cache.Set(token, token)
	if err != nil {
		return err
	}

	ttl := time.Until(claims.ExpiresAt.Time)
	if ttl < 0 {
		return nil
	}

	_, err = cache.SetTTL(token, ttl)
	if err != nil {
		return err
	}

	return nil
}

func SaveMedia(orgID string, filename string, header *multipart.FileHeader, size int64, c *gin.Context) (*models.SavedMedia, error) {
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(header, dst); err != nil {
		return nil, fmt.Errorf("Failed to save file: %v", err)
	}

	selfUrl := config.GetServerConfig().SELF_URL
	fileURL := fmt.Sprintf("%s/media/%s", selfUrl, filename)

	orgObjID, err := primitive.ObjectIDFromHex(orgID)
	if err != nil {
		return nil, fmt.Errorf("Invalid (orgID): %v", err)
	}

	saveMedia := models.SavedMedia{
		OrgID:    orgObjID,
		Filename: filename,
		URL:      fileURL,
		Size:     header.Size,
		SavedAt:  primitive.NewDateTimeFromTime(time.Now()),
	}

	err = repository.SaveMediaInRepo(&saveMedia)
	if err != nil {
		return nil, err
	}

	return &saveMedia, nil
}
