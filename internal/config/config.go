package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"lembrago.com/lembrago/database"
)

type AppConfig struct {
	JWTSecret             []byte
	JWTSecretUserCreation []byte
	JWTIssuer             string
	Host                  string
	Port                  string
}

func init() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Error ao conectar ao mongodb: ", err)
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
}

func GetAppConfig() *AppConfig {
	cfg := &AppConfig{
		JWTSecret:             []byte(os.Getenv("JWT_SECRET")),
		JWTSecretUserCreation: []byte(os.Getenv("JWT_SECRET_USER_CREATION")), // Carregue se usar
		JWTIssuer:             os.Getenv("JWT_ISSUER"),
		Host:                  os.Getenv("HOST"),
		Port:                  os.Getenv("PORT"),
	}

	return cfg
}
