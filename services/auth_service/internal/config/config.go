package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environment          string        `mapstructure:"ENVIRONMENT"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBAddress            string        `mapstructure:"DB_ADDRESS"`
	DBPort               int           `mapstructure:"DB_PORT"`
	DBUser               string        `mapstructure:"DB_USER"`
	DBPasswd             string        `mapstructure:"DB_PASSWD"`
	DBSName              string        `mapstructure:"DB_NAME"`
	RedisAddr            string        `mapstructure:"REDIS_ADDRESS"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	GRPCGWServerAddress  string        `mapstructure:"GRPC_GW_SERVER_ADDRESS"`
	AllowOrigins         string        `mapstructure:"HTTP_ALLOW_ORIGINS"`
	TokenSecretKey       string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`

	PROCESS_INTERVAL time.Duration `mapstructure:"PROCESS_INTERVAL"`
	DebugLv          int           `mapstructure:"DEBUG_LV"`
}

func LoadConfig(path string) (Config, error) {
	var config Config
	var err error = nil
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = viper.Unmarshal(&config)
	return config, nil
}
