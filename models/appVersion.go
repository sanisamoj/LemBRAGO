package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type AppVersion struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	ActualServerVersion string             `bson:"actualServerVersion" json:"actualServerVersion" validate:"required"`

	LatestDesktopVersion string `bson:"latestDesktopVersion" json:"latestDesktopVersion" validate:"required"`
	MinDesktopVersion    string `bson:"minDesktopVersion" json:"minDesktopVersion" validate:"required"`

	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}
