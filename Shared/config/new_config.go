package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Загрузчик параметров для получения конфига
type ConfigLoader struct {
	GlobalPathToEnv string //Дериктория, в которой находиться общий файл env
	EnvFile         string //Название файла с расщирением .env
	EnvType         string //Тип файла env
	ConfigType      string //Тип конфиг файла (yaml)
	PathToLocalEnv  string //Переменная окружения, в которой лежит путь в локально env файлу
	PathToConfig    string //Переменная окружения, в которой лежит путь к конфигу
}

// Конструктор для параметров
func NewConfigLoader(globalPathToEnv, envFile, configType, pathToLocalEnv, pathToConfig string) *ConfigLoader {
	loader := ConfigLoader{
		GlobalPathToEnv: globalPathToEnv,
		EnvFile:         envFile,
		EnvType:         "env",
		ConfigType:      configType,
		PathToLocalEnv:  pathToLocalEnv,
		PathToConfig:    pathToConfig,
	}
	return &loader
}

// Get a new config
func NewConfig[T any](loader *ConfigLoader) (*T, error) {

	pathLocalEnv, err := GetPathToEnv(loader)
	if err != nil {
		return nil, err
	}

	configViper, err := GetConfigViper(pathLocalEnv, loader)
	if err != nil {
		return nil, err
	}

	var cfg T
	if err := configViper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// Get path to local env
func GetPathToEnv(loader *ConfigLoader) (string, error) {
	viper := viper.New()
	viper.AddConfigPath(loader.GlobalPathToEnv)
	viper.SetConfigFile(loader.EnvFile)
	viper.SetConfigType(loader.EnvType)

	if err := viper.ReadInConfig(); err != nil {
		return "", err
	}

	pathToLocalEnv := viper.GetString(loader.PathToLocalEnv)
	if pathToLocalEnv == "" {
		return "", fmt.Errorf("variable %s is not set", loader.PathToLocalEnv)
	}
	return pathToLocalEnv, nil
}

func GetConfigViper(pathLocalEnv string, loader *ConfigLoader) (*viper.Viper, error) {
	envViper := viper.New()
	envViper.SetConfigFile(pathLocalEnv)
	envViper.SetConfigType(loader.EnvType)

	if err := envViper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading .env: %w", err)
	}

	pathConfig := envViper.GetString(loader.PathToConfig)
	if pathConfig == "" {
		return nil, fmt.Errorf("variable %s is not set", loader.PathToConfig)
	}

	configViper := viper.New()
	configViper.SetConfigFile(pathConfig)
	configViper.SetConfigType(loader.ConfigType)

	if err := configViper.ReadInConfig(); err != nil {
		return nil, err
	}

	for _, key := range envViper.AllKeys() {
		configKey := strings.ToLower(strings.ReplaceAll(key, "_", "."))
		if val := envViper.Get(key); val != nil {
			configViper.Set(configKey, val)
		}
	}
	return configViper, nil
}
