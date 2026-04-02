package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig      `mapstructure:"app"`
	JWT      JWTConfig      `mapstructure:"jwt"`
	Database DatabaseConfig `mapstructure:"database"`
	Valkey   ValkeyConfig   `mapstructure:"valkey"`
	Agora    AgoraConfig    `mapstructure:"agora"`
	Google   GoogleConfig   `mapstructure:"google"`
}

type AppConfig struct {
	Port              int    `mapstructure:"port"`
	Env               string `mapstructure:"env"`
	MatchmakingEngine string `mapstructure:"matchmaking_engine"` // "valkey" or "postgres"
}

type JWTConfig struct {
	Secret     string `mapstructure:"secret"`
	Expiration int    `mapstructure:"expiration"`
}

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"ssl_mode"`
}

type ValkeyConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

type AgoraConfig struct {
	AppID          string `mapstructure:"app_id"`
	AppCertificate string `mapstructure:"app_certificate"`
}

type GoogleConfig struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	RedirectURL  string `mapstructure:"redirect_url"`
	GeminiAPIKey string `mapstructure:"gemini_api_key"`
}

func LoadConfig(path string) (*Config, error) {
	viper.AddConfigPath(path)
	viper.AddConfigPath(".")
	viper.AddConfigPath("../../")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	log.Println("Configuration Load Success")

	return &config, nil
}
