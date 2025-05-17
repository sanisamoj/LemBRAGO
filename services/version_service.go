package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"lembrago.com/lembrago/internal/config"
	"lembrago.com/lembrago/models"
	"lembrago.com/lembrago/repository"
)

func RegisterVersion(req *models.ApplicationVersion) (*models.ApplicationVersion, error) {
	req.ID = primitive.NewObjectID()
	req.CreatedAt = primitive.NewDateTimeFromTime(time.Now())
	req.LatestDesktopVersion.PubDate = primitive.NewDateTimeFromTime(time.Now())

	err := repository.RegisterAppVersion(req)
	if err != nil {
		return nil, err
	}
	req.ActualServerVersion = config.GetServerVersion()
	return req, nil
}

func GetLatestVersion() (*models.ApplicationVersion, error) {
	version, err := repository.GetLastestVersion()
	if err != nil {
		return nil, err
	}

	version.ActualServerVersion = config.GetServerVersion()
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

	for i := range allVersions {
		allVersions[i].ActualServerVersion = config.GetServerVersion()
	}

	return allVersions, nil
}

func UpdateVersion(req *models.ApplicationVersion) (*models.ApplicationVersion, error) {
	return repository.UpdateAppVersion(req)
}
