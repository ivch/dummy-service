package config

import (
	"os"

	"gopkg.in/go-playground/validator.v9"
)

type Config struct {
	HTTPPort     string `validate:"required"`
	HTTPSPort    string `validate:"required_without=HTTPPort"`
	AWSConfig    AWS
	DynamoConfig Dynamo
}

type AWS struct {
	Profile string
	Key     string `validate:"required"`
	Secret  string `validate:"required"`
	Region  string `validate:"required"`
}

type Dynamo struct {
	Host  string
	Table string `validate:"required"`
}

func New() (*Config, error) {
	cfg := Config{
		HTTPPort:  os.Getenv("HTTP_PORT"),
		HTTPSPort: os.Getenv("HTTPS_PORT"),
		AWSConfig: AWS{
			Profile: os.Getenv("AWS_PROFILE"),
			Key:     os.Getenv("AWS_KEY"),
			Secret:  os.Getenv("AWS_SECRET"),
			Region:  os.Getenv("AWS_REGION"),
		},
		DynamoConfig: Dynamo{
			Host:  os.Getenv("DB_HOST"),
			Table: os.Getenv("DB_TABLE"),
		},
	}

	if err := validator.New().Struct(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
