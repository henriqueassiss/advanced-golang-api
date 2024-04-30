package config

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	Api
	App
	Client
	Cors

	Cache
	Database
}

func New(isTesting, isCallingFromDomain bool) *Config {
	var err error

	if isTesting {
		LoadEnvTest(isCallingFromDomain)
	} else {
		err = godotenv.Load()
	}

	if err != nil {
		log.Println(err)
	}

	if isTesting {
		return &Config{
			Api:      API(),
			App:      APP(),
			Cache:    NewCache(),
			Database: DataStore(),
		}
	}

	return &Config{
		Api:      API(),
		App:      APP(),
		Client:   NewClient(),
		Cors:     NewCors(),
		Cache:    NewCache(),
		Database: DataStore(),
	}
}

func LoadEnvTest(isCallingFromDomain bool) {
	rootDir, err := "", error(nil)
	if isCallingFromDomain {
		rootDir, err = filepath.Abs("../../..")
	} else {
		rootDir, err = filepath.Abs("../..")
	}

	if err != nil {
		log.Fatalf("Error getting project root directory: %v", err)
	}

	envFilePath := filepath.Join(rootDir, ".env.test")
	err = godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env.test file: %v", err)
	}
}
