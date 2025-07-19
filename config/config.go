package config

import (
	"github.com/joho/godotenv"

	kenv "github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

type (
	Config struct {
		SmtpHost         string `koanf:"SMTP_HOST"`
		SmtpPort         int    `koanf:"SMTP_PORT"`
		SmtpSenderName   string `koanf:"SMTP_SENDER_NAME"`
		SmtpAuthEmail    string `koanf:"SMTP_AUTH_EMAIL"`
		SmtpAuthPassword string `koanf:"SMTP_AUTH_PASSWORD"`
	}
)

func LoadConfig() (*Config, error) {
	var config Config

	_ = godotenv.Load(".env")

	k := koanf.New(".")

	k.Load(kenv.Provider("", ".", func(s string) string {
		return s
	}), nil)

	if err := k.Unmarshal("", &config); err != nil {
		return nil, err
	}

	return &config, nil
}
