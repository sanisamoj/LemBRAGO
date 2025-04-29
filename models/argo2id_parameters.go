package models

type Argo2IDParameters struct {
	Memory      uint32 `json:"memory"`
	Time        uint32 `json:"time"`
	Parallelism uint8  `json:"parallelism"`
	SaltLength  uint32 `json:"saltLength"`
	KeyLength   uint32 `json:"keyLength"`
}
