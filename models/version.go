package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ApplicationVersion struct {
	ID                  primitive.ObjectID `bson:"_id" json:"id"`
	Version   string              `bson:"version" json:"version" validate:"required"`
	Notes     string              `bson:"notes" json:"notes" validate:"required"`
	PubDate   primitive.DateTime  `bson:"pubDate" json:"pub_date"`
	Platforms map[string]Platform `bson:"platforms" json:"platforms" validate:"required"`
	Type      string              `bson:"type" json:"type" validate:"required"`
	Changes   []string            `bson:"changes" json:"changes" validate:"required"`
	CreatedAt primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

type Platform struct {
	Signature string `bson:"signature" json:"signature" validate:"required"`
	Url       string `bson:"url" json:"url" validate:"required"`
}
