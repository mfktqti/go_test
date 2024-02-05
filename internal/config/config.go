package config

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	localConfigFile string
	localConfig     Config
)

func init() {
	flag.StringVar(&localConfigFile, "c", "./config/config.yaml", "init local config")
	// flag.StringVar(&localConfigFile, "c", "../../../config/config.yaml", "init local config")
}

// LoadConfig 加载本地配置文件
func LoadConfig() (*Config, error) {
	configBytes, err := os.ReadFile(localConfigFile)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %w", err)
	}

	if err = yaml.Unmarshal(configBytes, &localConfig); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w, file: %v", err, localConfigFile)
	}

	return &localConfig, nil
}

// GetConfig 获取配置
func GetConfig() *Config {
	return &localConfig
}

// LogConfig 配置项
type LogConfig struct {
	LogPath    string `yaml:"logPath"`
	LogLevel   int8   `yaml:"logLevel"`
	MaxSize    int    `yaml:"maxSize"`
	MaxAge     int    `yaml:"maxAge"`
	MaxBackups int    `yaml:"maxBackups"`
	OutputFile bool   `yaml:"outputFile"`
	StackTrace bool   `yaml:"stackTrace"`
	Compress   bool   `yaml:"compress"`
}

// Config 配置
type Config struct {

	// Etcd  *EtcdConfig  `yaml:"etcd"`
	Mysql *MysqlConfig `yaml:"mysql"`
	// Redis *RedisConfig `yaml:"redis"`
	Log *LogConfig `yaml:"log"`
}

// MysqlConfig 配置项
type MysqlConfig struct {
	Addr         string `yaml:"path"`
	UserName     string `yaml:"username"`
	Password     string `yaml:"password"`
	DbName       string `yaml:"db-name"`
	MaxOpenConns int    `yaml:"max-open-conns"`
	MaxIdleConns int    `yaml:"max-idle-conns"`
}

// RedisConfig 配置项
type RedisConfig struct {
	Addr        string `yaml:"addr"`
	Db          int    `yaml:"db"`
	Password    string `yaml:"password"`
	MaxActive   int    `yaml:"maxActive"`
	MaxIdle     int    `yaml:"maxIdle"`
	IdleTimeout int    `yaml:"idleTimeout"`
}

// EtcdConfig 配置项
type EtcdConfig struct {
	Addr string `yaml:"addr"`
}
