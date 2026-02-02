package config

import (
	"time"

	"github.com/spf13/viper"
)

// RunMode 실행 모드
type RunMode int

const (
	ModeDev RunMode = 0 // 개발 모드
	ModeOpr RunMode = 1 // 운영 모드
)

func (m RunMode) String() string {
	switch m {
	case ModeDev:
		return "development"
	case ModeOpr:
		return "production"
	default:
		return "unknown"
	}
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

	AUTH_SERVICE_URL_DEV   string `mapstructure:"AUTH_SERVICE_URL_DEV"`
	AUTH_SERVICE_URL_OPR   string `mapstructure:"AUTH_SERVICE_URL_OPR"`
	DOCKER_SERVICE_URL_DEV string `mapstructure:"DOCKER_SERVICE_URL_DEV"`
	DOCKER_SERVICE_URL_OPR string `mapstructure:"DOCKER_SERVICE_URL_OPR"`

	AUTH_SERVICE_URL   string
	DOCKER_SERVICE_URL string

	Mode RunMode // 현재 실행 모드
}

func LoadConfig(path string, mode RunMode) (Config, error) {
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
	if err != nil {
		return config, err
	}

	// 모드에 따라 서비스 URL 설정
	config.Mode = mode
	config.SetServiceURLs()

	return config, nil
}

// SetServiceURLs 모드에 따라 서비스 URL 설정
func (c *Config) SetServiceURLs() {
	switch c.Mode {
	case ModeOpr:
		c.AUTH_SERVICE_URL = c.AUTH_SERVICE_URL_OPR
		c.DOCKER_SERVICE_URL = c.DOCKER_SERVICE_URL_OPR
	default: // ModeDev
		c.AUTH_SERVICE_URL = c.AUTH_SERVICE_URL_DEV
		c.DOCKER_SERVICE_URL = c.DOCKER_SERVICE_URL_DEV
	}
}
