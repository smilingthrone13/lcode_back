package config

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"lcode/pkg/digit"
	"os"
	"time"
)

const (
	defaultHTTPRWTimeout          = 10 * time.Second
	defaultHTTPMaxHeaderMegabytes = 1

	defaultPage       = 1
	defaultLimitCount = 30

	defaultUserAvatarMaxSize = 5 * 1024 * 1024 // MB

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
		Files             Files
		DBConfig          DBConfig
		QueryParams       QueryParams
		JudgeConfig       JudgeConfig
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

	Files struct {
		MainFolder        string
		UserAvatarMaxSize int64
	}

	JudgeConfig struct {
		Host                 string
		Port                 string
		DefaultMemoryLimitKB int     `mapstructure:"defaultMemoryLimitKB"`
		DefaultTimeLimitSec  float64 `mapstructure:"defaultTimeLimitSec"`
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

func parseFiles(cfg *Config) error {
	var f struct {
		MainFolder        string
		UserAvatarMaxSize string
	}

	if err := viper.UnmarshalKey("files.mainFolder", &f.MainFolder); err != nil {
		return err
	}

	cfg.Files.MainFolder = f.MainFolder

	if err := viper.UnmarshalKey("files.userAvatarMaxSize", &f.UserAvatarMaxSize); err != nil {
		return err
	}

	size, err := digit.ParseSize(f.UserAvatarMaxSize)
	if err != nil {
		return err
	}

	cfg.Files.UserAvatarMaxSize = size

	return nil
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

	if err := parseFiles(cfg); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("tls", &cfg.TLS); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("query_params", &cfg.QueryParams); err != nil {
		return err
	}

	if err := viper.UnmarshalKey("judge", &cfg.JudgeConfig); err != nil {
		return err
	}

	return nil
}

func parseEnv(configDir string, cfg *Config) error {
	path_db, ok := os.LookupEnv("PATH_DB")
	if ok {
		cfg.DBConfig.Path = path_db

		return nil
	}

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
	viper.SetDefault("query_params.limit", defaultLimitCount)

	viper.SetDefault("files.user_avatar_max_size", defaultUserAvatarMaxSize)

	viper.SetDefault("auth.access_token_exp_time", defaultAccessTokenExpTime)
	viper.SetDefault("auth.refresh_token_exp_time", defaultRefreshTokenExpTime)
	viper.SetDefault("auth.secret", defaultSecretKey)
}
