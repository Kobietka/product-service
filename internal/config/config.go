package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseUrl string
	Port        string
}

func NewConfigStore() Store {
	return Store{}
}

type Store struct{}

func (s Store) GetConfig() (Config, error) {
	databaseUrl, ok := os.LookupEnv("DATABASE_URL")
	if !ok {
		return Config{}, errors.New("DATABASE_URL environment variable not set")
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		return Config{}, errors.New("PORT environment variable not set")
	}

	return Config{
		DatabaseUrl: databaseUrl,
		Port:        port,
	}, nil
}
