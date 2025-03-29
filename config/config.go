package config

import (
	"os"
	"time"

	
)

type Config struct {
	App         string
	Environment string
	LogLevel    string
	Server      struct {
		Host        string
		Port         string
		ReadTimeout  string
		WriteTimeout string
		IdleTimeout  string
	}

	Context struct {
		Timeout time.Duration
	}

	DB struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SslMode  string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		Name     string
		Time     time.Time
	}
	Token struct {
		Secret     string
		AccessTTL  time.Duration
		RefreshTTL time.Duration
		SignInKey  string
	}
	Minio struct {
		Endpoint                 string
		AccessKeyID              string
		SecretAcessKey           string
		Location                 string
		ImageUrlUploadBucketName string
		FileUploadBucketName     string
	}
	SMTP struct {
		Email         string
		EmailPassword string
		SMTPPort      string
		SMTPHost      string
	}

}

func NewConfig() (*Config, error) {
	var config Config

	// general configuration
	config.App = getEnv("APP", "app")
	config.Environment = getEnv("ENVIRONMENT", "develop")
	config.LogLevel = getEnv("LOG_LEVEL", "debug")

	// server configuration
	config.Server.Host = getEnv("SERVER_HOST", "localhost")  //app
	config.Server.Port = getEnv("SERVER_PORT", ":7777")
	config.Server.ReadTimeout = getEnv("SERVER_READ_TIMEOUT", "10s")
	config.Server.WriteTimeout = getEnv("SERVER_WRITE_TIMEOUT", "10s")
	config.Server.IdleTimeout = getEnv("SERVER_IDLE_TIMEOUT", "120s")

	//context configuration
	ContexTimeout, err := time.ParseDuration(getEnv("CONTEXT_TIMEOUT", "30s"))
	if err != nil {
		return nil, err
	}

	config.Context.Timeout = ContexTimeout

	// db configuration
	config.DB.Host = getEnv("POSTGRES_HOST", "localhost") // postgres
	config.DB.Port = getEnv("POSTGRES_PORT", "5432")
	config.DB.User = getEnv("POSTGRES_USER", "postgres")
	config.DB.Password = getEnv("POSTGRES_PASSWORD", "")
	config.DB.SslMode = getEnv("POSTGRES_SSLMODE", "disable")
	config.DB.Name = getEnv("POSTGRES_DATABASE", "hospitaldb")

	// access ttl parse
	accessTTl, err := time.ParseDuration(getEnv("TOKEN_ACCESS_TTL", "3h"))
	if err != nil {
		return nil, err
	}
	// refresh ttl parse
	refreshTTL, err := time.ParseDuration(getEnv("TOKEN_REFRESH_TTL", "24h"))
	if err != nil {
		return nil, err
	}
	config.Token.AccessTTL = accessTTl
	config.Token.RefreshTTL = refreshTTL
	config.Token.SignInKey = getEnv("TOKEN_SIGNIN_KEY", "debug")

	// redis configuration
	config.Redis.Host = getEnv("REDIS_HOST", "localhost")  //redisdb
	config.Redis.Port = getEnv("REDIS_PORT", "6379")
	config.Redis.Password = getEnv("REDIS_PASSWORD", "")
	config.Redis.Name = getEnv("REDIS_DATABASE", "0")

	//smtp confifuration
	config.SMTP.Email = getEnv("SMTP_EMAIL", "theuniver77@gmail.com")
	config.SMTP.EmailPassword = getEnv("SMTP_EMAIL_PASSWORD", "")
	config.SMTP.SMTPPort = getEnv("SMTP_PORT", "587")
	config.SMTP.SMTPHost = getEnv("SMTP_HOST", "smtp.gmail.com")

	return &config, nil
}


func getEnv(key string, defaultVaule string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	}
	return defaultVaule
}
