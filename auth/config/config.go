package config

import (
	"errors"
	"log"
	"time"

	"github.com/spf13/viper"
)

// App config struct
type Config struct {
	Server       ServerConfig
	Postgres     PostgresConfig
	Redis        RedisConfig
	Kafka        KafkaConfig
	RefreshToken RefreshToken
	AccessToken  AccessToken
	Logger       Logger
	SMTP         SMTP
	Gateway      Gateway
	Casbin       Casbin
	Metrics      Metrics
	Jaeger       Jaeger
}

// Server config struct
type ServerConfig struct {
	AppVersion        string
	Port              string
	Mode              string
	JwtSecretKey      string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	Debug             bool
	MaxConnectionIdle time.Duration
	Timeout           time.Duration
	MaxConnectionAge  time.Duration
	Time              time.Duration
}

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}

// Postgres config
type PostgresConfig struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDbname   string
	PostgresSSLMode  string
}

// Redis config
type RedisConfig struct {
	RedisAddr      string
	RedisPassword  string
	RedisDB        string
	RedisDefaultdb string
	MinIdleConns   int
	PoolSize       int
	PoolTimeout    int
	Password       string
	DB             int
}

type KafkaConfig struct {
	Brokers []string
	GroupID string
}

// RefreshToken config
type RefreshToken struct {
	MaxAge time.Duration
}

// AccessToken config
type AccessToken struct {
	Expire time.Duration
}

// SMTP config
type SMTP struct {
	Host           string
	Port           int
	User           string
	Pass           string
	SenderIdentity string
	SenderEmail    string
}

// Server config struct
type Gateway struct {
	Port string
}

// Casbin Config
type Casbin struct {
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresSSLMode  string
	URL              string
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Jaeger
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}

// Load config file from given path
func LoadConfig(filename string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigName(filename)
	v.AddConfigPath(".")
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("config file not found")
		}
		return nil, err
	}

	return v, nil
}

// Parse config file
func ParseConfig(v *viper.Viper) (*Config, error) {
	c := new(Config)

	err := v.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return c, nil
}

// Get config
func GetConfig(configPath string) (*Config, error) {
	cfgFile, err := LoadConfig(configPath)
	if err != nil {
		return nil, err
	}

	cfg, err := ParseConfig(cfgFile)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}
