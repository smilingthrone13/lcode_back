package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

const (
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	defaultPage       = 1
	defaultLimitCount = 30

	defaultAccessTokenExpTime  = time.Second * 300
	defaultRefreshTokenExpTime = time.Hour * 24 * 30
	defaultSecretKey           = "secret"
)

type (
	Config struct {
		IsDebug           bool
		CorsOrigins       []string
		SearchCoefficient float32
		HTTP              HTTPConfig
		Auth              AuthConfig
		TLS               TLSConfig
		DBConfig          DBConfig
		QueryParams       QueryParams
	}

	HTTPConfig struct {
		Host               string
		Port               string
		ReadTimeout        time.Duration
		WriteTimeout       time.Duration
		MaxHeaderMegabytes int
	}

	TLSConfig struct {
		Enabled  bool
		CertFile string `mapstructure:"cert"`
		KeyFile  string `mapstructure:"key"`
	}

	AuthConfig struct {
		AccessTokenExpTime  time.Duration
		RefreshTokenExpTime time.Duration
		Secret              string
	}

	DBConfig struct {
		Path string
	}

	QueryParams struct {
		Limit int
		Page  int
	}
)

func Init(configsDir string) (*Config, error) {
	InitDefault()
	var cfg Config

	err := parseYml(configsDir, &cfg)
	if err != nil {
		return nil, err
	}

	err = parseEnv(configsDir, &cfg)
	if err != nil {
		return nil, err
	}

	viper.Reset()

	return &cfg, nil
}

func parseYml(configDir string, cfg *Config) error {
	if err := parseConfigFile(configDir+"config", "yaml"); err != nil {
		fmt.Print(err.Error())

		return err
	}

	if err := viper.UnmarshalKey("isDebug", &cfg.IsDebug); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("corsOrigins", &cfg.CorsOrigins); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("searchCoefficient", &cfg.SearchCoefficient); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("http", &cfg.HTTP); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("auth", &cfg.Auth); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("tls", &cfg.TLS); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("query_params", &cfg.QueryParams); err != nil {
		return err
	}

	return nil
}

func parseEnv(configDir string, cfg *Config) error {
	if err := parseConfigFile(configDir, "env"); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("PATH_DB", &cfg.DBConfig.Path); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(folder string, fileType string) error {
	viper.AddConfigPath(folder)
	switch fileType {
	case "yaml":
		viper.SetConfigName("config")
	case "env":
		viper.SetConfigName("App")
		viper.AutomaticEnv()
	default:
		return errors.New("fileType is invalid")
	}

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.MergeInConfig()
}

func InitDefault() {
	viper.SetDefault("http.max_header_megabytes", defaultHTTPMaxHeaderMegabytes)
	viper.SetDefault("http.timeouts.read", defaultHTTPRWTimeout)
	viper.SetDefault("http.timeouts.write", defaultHTTPRWTimeout)
	viper.SetDefault("query_params.page", defaultPage)

	viper.SetDefault("auth.access_token_exp_time", defaultAccessTokenExpTime)
	viper.SetDefault("auth.refresh_token_exp_time", defaultRefreshTokenExpTime)
	viper.SetDefault("auth.secret", defaultSecretKey)
}
