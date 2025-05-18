package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/repository"
)

func RegisterVersion(req *models.ApplicationVersion) (*models.ApplicationVersion, error) {
	req.ID = primitive.NewObjectID()
	req.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	req.PubDate = primitive.NewDateTimeFromTime(time.Now())

	err := repository.RegisterAppVersion(req)
	if err != nil {
		return nil, err
	}
	return req, nil
}

func GetLatestVersion() (*models.ApplicationVersion, error) {
	version, err := repository.GetLastestVersion()
	if err != nil {
		return nil, err
	}

	return version, nil
}

func GetAllVersions() ([]models.ApplicationVersion, error) {
	allVersions, err := repository.GetAllVersions()
	if err != nil {
		return nil, err
	}
	if allVersions == nil {
		return []models.ApplicationVersion{}, nil
	}

	return allVersions, nil
}

func UpdateVersion(req *models.ApplicationVersion) (*models.ApplicationVersion, error) {
	return repository.UpdateAppVersion(req)
}
