package config

import (
	"os"
)

type Config struct {
	Port     string
	AuthUser string
	AuthPass string
	DB       string
	DB_DIR   string
}

const (
	DEFAULT_PORT      = "8080"
	DEFAULT_AUTH_USER = "heroku"
	DEFAULT_DB        = "badger"
	DEFAULT_DB_DIR    = "/tmp/plivo/contact-book"
	DEFAULT_AUTH_PASS = "plivo"
)

func GetConfig() *Config {
	c := Config{}
	c.Port = os.Getenv("PORT")
	if c.Port == "" {
		c.Port = DEFAULT_PORT
	}

	c.AuthPass = os.Getenv("PLIVO_AUTH_PASS")
	if c.AuthPass == "" {
		c.AuthPass = DEFAULT_AUTH_PASS
	}

	c.AuthUser = os.Getenv("PLIVO_AUTH_USER")
	if c.AuthUser == "" {
		c.AuthUser = DEFAULT_AUTH_USER
	}

	c.DB = os.Getenv("PLIVO_CONTACTS_DB")
	if c.DB == "" {
		c.DB = DEFAULT_DB
	}

	c.DB_DIR = os.Getenv("PLIVO_DB_DIR")
	if c.DB_DIR == "" {
		c.DB_DIR = DEFAULT_DB_DIR
	}

	return &c
}
