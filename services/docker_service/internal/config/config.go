package config

import (
	"encoding/json"
	"time"

	"github.com/spf13/viper"
)

// DockerHostConfig는 Docker 호스트 설정
type DockerHostConfig struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
}

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
	AllowOrigins         string        `mapstructure:"HTTP_ALLOW_ORIGINS"`
	TokenSecretKey       string        `mapstructure:"TOKEN_SYMMETRIC_KEY"`

	PROCESS_INTERVAL time.Duration `mapstructure:"PROCESS_INTERVAL"`
	DebugLv          int           `mapstructure:"DEBUG_LV"`

	CERT_PATH    string `mapstructure:"CERT_PATH"`
	DOCKER_HOSTS string `mapstructure:"DOCKER_HOSTS"` // JSON format: [{"name":"host1","addr":"tcp://..."}]
}

// GetDockerHosts는 DOCKER_HOSTS JSON 문자열을 파싱하여 반환
func (c *Config) GetDockerHosts() ([]DockerHostConfig, error) {
	if c.DOCKER_HOSTS == "" {
		return []DockerHostConfig{}, nil
	}

	var hosts []DockerHostConfig
	if err := json.Unmarshal([]byte(c.DOCKER_HOSTS), &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
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
