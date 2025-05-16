package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/repository"
)

func RegisterVersion(req *models.AppVersion) (*models.AppVersion, error) {
	req.ID = primitive.NewObjectID()
	req.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	err := repository.RegisterAppVersion(req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func GetLatestVersion() (*models.AppVersion, error) {
	return repository.GetLastestAppVersion()
}

func GetAllVersions() ([]models.AppVersion, error) {
	allVersions, err := repository.GetAllAppVersion()
	if err != nil {
		return nil, err
	}
	if allVersions == nil  {
		return []models.AppVersion{}, nil
	}

	return allVersions, nil
}

func UpdateVersion(req *models.AppVersion) (*models.AppVersion, error) {
	return repository.UpdateAppVersion(req)
}