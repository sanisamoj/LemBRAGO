package models

type KeysDTO struct {
	PublicKey           string          `json:"publicKey" validate:"required"`
	EncryptedPrivateKey EncryptedKeyDto `json:"encryptedPrivateKey" validate:"required"` // EUserPrivK
	EncryptedSecretKey  EncryptedKeyDto `json:"encryptedSecretKey" validate:"required"`  // ESK
}

type EncryptedKeyDto struct {
	Ciphertext string `json:"ciphertext" validate:"required"`
	Nonce      string `json:"nonce" validate:"required"`
}
