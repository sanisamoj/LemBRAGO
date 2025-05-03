package models

type CreateUserRequest struct {
	Code     string `json:"code"`
	Username string `bson:"username" json:"username" validate:"required"`

	PasswordVerifier struct { // PV, salt_pv and paramaters
		Salt       string            `json:"salt" validate:"required"`
		Verifier   string            `json:"verifier" validate:"required"`
		Parameters Argo2IDParameters `json:"parameters" validate:"required"`
	} `json:"passwordVerifier"`

	Salt_ek string `json:"salt_ek" validate:"required"`

	Keys KeysDTO `json:"keys"`

	MyVault *CreateVaultRequest `json:"myVault"` // `json:"myVault" validate:"required"`
}

type UserResponse struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	OrgId    string   `json:"orgId"`
	Username string   `bson:"username" json:"username" validate:"required"`
	Role     UserRole `json:"role"`

	PasswordVerifier PasswordVerifierResponse `json:"passwordVerifier"`
	Salt_ek          string                   `json:"salt_ek" validate:"required"`

	Keys KeysDTO `json:"keys"`
}

type UserWithOrganizationResponse struct {
	OrgID            string                   `json:"orgId"`
	OrganizationName string                   `json:"organizationName"`
	OrgImagUrl       string                   `json:"organizationImageUrl"`
	PasswordVerifier PasswordVerifierResponse `json:"passwordVerifier"`
}

type PasswordVerifierResponse struct {
	Salt       string            `json:"salt" validate:"required"`
	Parameters Argo2IDParameters `json:"parameters" validate:"required"`
}

type UserLoginResponse struct {
	User  UserResponse `json:"user" validate:"required"`
	Token string       `json:"token" validate:"required"`
}

type UserLoginComparison struct {
	OrgID    string `json:"orgId" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Verifier string `json:"verifier" validate:"required"`
}

type InviteUserRequest struct {
	Email string   `json:"email" validate:"required,email"`
	Role  UserRole `json:"role" validate:"required,oneof=admin member"`
}

type MinimalUserInfoResponse struct {
	ID        string     `json:"id"`
	OrgID     string     `json:"orgId"`
	Email     string     `json:"email"`
	Username  string     `json:"username"`
	PublicKey string     `json:"publicKey"`
	Role      UserRole   `json:"role"`
	Status    UserStatus `json:"status"`
	CreatedAt string     `json:"createdAt"`
	UpdatedAt string     `json:"updatedAt"`
}

type AuthCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type AuthCodeSendRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required"`
}
