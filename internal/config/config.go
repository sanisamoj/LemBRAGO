package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"lembrago.com/lembrago/database"
)

type ServerConfig struct {
	JWTSecret             []byte
	JWTSecretUserCreation []byte
	JWTSecretAdmin        []byte
	JWTIssuer             string
	Host                  string
	Port                  string
	SELF_URL              string
}

func init() {
	if err := database.ConnectDB(); err != nil {
		log.Fatal("Error ao conectar ao mongodb: ", err)
	}

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll("uploads", os.ModePerm)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat("releases"); os.IsNotExist(err) {
		_ = os.MkdirAll("releases, os.ModePerm", os.ModePerm)

		_ = os.MkdirAll("releases/windows-x86_64", os.ModePerm)
		_ = os.MkdirAll("releases/linux-x86_64", os.ModePerm)
	}

	// go IconPopulate()
}

func GetServerConfig() *ServerConfig {
	cfg := &ServerConfig{
		JWTSecret:             []byte(os.Getenv("JWT_SECRET")),
		JWTSecretUserCreation: []byte(os.Getenv("JWT_SECRET_USER_CREATION")),
		JWTSecretAdmin:        []byte(os.Getenv("JWT_SECRET_ADMIN")),
		JWTIssuer:             os.Getenv("JWT_ISSUER"),
		Host:                  os.Getenv("HOST"),
		Port:                  os.Getenv("PORT"),
		SELF_URL:              os.Getenv("SELF_URL"),
	}

	return cfg
}

func GetServerVersion() string {
	return "0.8.0"
}
