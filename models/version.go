package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ApplicationVersion struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	ActualServerVersion string             `bson:"actualServerVersion,omitempty" json:"actualServerVersion"`

	LatestDesktopVersion Version `bson:"latestDesktopVersion" json:"latestDesktopVersion" validate:"required"`
	MinDesktopVersion    Version `bson:"minDesktopVersion" json:"minDesktopVersion" validate:"required"`

	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

type Version struct {
	Version   string              `bson:"version" json:"version" validate:"required"`
	Notes     string              `bson:"notes" json:"notes" validate:"required"`
	PubDate   primitive.DateTime  `bson:"pubDate" json:"pub_date"`
	Platforms map[string]Platform `bson:"platforms" json:"platforms" validate:"required"`
}

type Platform struct {
	Signature string `bson:"signature" json:"signature" validate:"required"`
	Url       string `bson:"url" json:"url" validate:"required"`
}
