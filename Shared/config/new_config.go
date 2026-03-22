package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
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

// Создание нового конфига
func NewConfig[T any](loader *ConfigLoader) (*T, error) {

	// Глобальный Viper. Читает переменные окружения из внешнего env файла, в котором указаны пути к локальным файлам env в сервисах
	globalViper := viper.New()
	globalViper.AddConfigPath(loader.GlobalPathToEnv)
	globalViper.SetConfigFile(loader.EnvFile)
	globalViper.SetConfigType(loader.EnvType)

	//Прочтения содержимого файла
	err := globalViper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	//Использование переменных окружения и получения пути к env файлу сервиса
	globalViper.AutomaticEnv()

	pathLocalEnv := globalViper.GetString(loader.PathToLocalEnv)
	// fmt.Println(pathLocalEnv)

	//Создание локального viper для работы внутри сервиса. Поиск файла по полученному адресу
	localViper := viper.New()
	localViper.SetConfigFile(pathLocalEnv)
	localViper.SetConfigType(loader.EnvType)

	//Прочтение содержимого env файла
	err = localViper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения env файла из %s: %w", pathLocalEnv, err)
	}
	// fmt.Println("Все ключи из localViper:", localViper.AllKeys())

	//Использование переменных окружения и получения пути к конфигу сервиса
	localViper.AutomaticEnv()
	pathConfig := localViper.GetString(loader.PathToConfig)
	// fmt.Println(pathConfig)

	//Подгрузка переменных окружения из env файла сервиса
	err = godotenv.Load(pathLocalEnv)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки локальных переменных окружения:%w", err)
	}

	//Получение данных из конфига
	content, err := os.ReadFile(pathConfig)
	if err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг:%s.%w", pathConfig, err)
	}

	//Применение значений из переменных окружения
	expandedContent := os.ExpandEnv(string(content))
	// fmt.Println(expandedContent)

	//Поиск файла по полученному адресу
	localViper.SetConfigFile(pathConfig)
	localViper.SetConfigType(loader.ConfigType)

	//Прочтение заполненного конфига
	err = localViper.ReadConfig(strings.NewReader(expandedContent))
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			return nil, fmt.Errorf("конфиг файл не найден:%w", err)
		default:
			return nil, fmt.Errorf("ошибка прочтения конфига:%w", err)
		}
	}

	//Парсиг конфига
	var cfg T
	if err := localViper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("ошибка парсинга конфига: %w", err)
	}

	return &cfg, nil
}
