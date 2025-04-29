package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Organization struct {
	ID                 primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name               string             `bson:"name" json:"name" validate:"required"`
	Email              string             `bson:"email" json:"email" validate:"required,email"`
	ImageUrl           string             `bson:"imageUrl,omitempty" json:"imageUrl"`
	SubscriptionPlan   SubscriptionPlan   `bson:"subscriptionPlan" json:"subscriptionPlan" validate:"required"`
	SubscriptionStatus string             `bson:"subscriptionStatus" json:"subscriptionStatus"`
	UpdatedAt          primitive.DateTime `bson:"updatedAt" json:"updatedAt"`
	CreatedAt          primitive.DateTime `bson:"createdAt" json:"createdAt"`
}

type CreateOrganizationRequest struct {
	Name             string            `json:"name" validate:"required"`
	Email            string            `json:"email" validate:"required,email"`
	SubscriptionPlan string            `json:"subscriptionPlan" validate:"required"`
	User             CreateUserRequest `json:"user" validate:"required"`
	ImageUrl         string            `json:"imageUrl,omitempty"`
}

type MinOrgWithTokenResponse struct {
	Token  string `json:"token"`
	Organization   string `json:"organization"`
	ImgUrl string `json:"imageUrl"`
	UserEmail string `json:"userEmail"`
}

type SubscriptionPlan string
type SubscriptionStatus string

const (
	UniquePlan SubscriptionPlan = "unique"
	BasicPlan  SubscriptionPlan = "basic"

	TrialStatus    SubscriptionStatus = "trial"
	ActiveStatus   SubscriptionStatus = "active"
	CanceledStatus SubscriptionStatus = "canceled"
	ExpiredStatus  SubscriptionStatus = "expired"
)
