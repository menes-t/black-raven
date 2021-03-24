package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/menes-t/black-raven/config/model"
	"github.com/menes-t/black-raven/logger"
	"github.com/spf13/viper"
	"sync"
)

type Config struct {
	viperInstance *viper.Viper
	mutex         *sync.RWMutex
	config        *ApplicationConfig
}

type ApplicationConfig struct {
	Tasks        []model.TaskConfig
	StartingTime string
}

type ApplicationConfigGetter interface {
	GetConfig() ApplicationConfig
}

func NewApplicationConfigGetter(path string) (ApplicationConfigGetter, error) {

	viperInstance := getViperInstance(path)
	viperInstance.SetTypeByDefaultValue(true)
	err := viperInstance.ReadInConfig()
	if err != nil {
		return nil, errors.New("config could not be found")
	}
	viperInstance.WatchConfig()
	configuration := &Config{viperInstance, &sync.RWMutex{}, nil}

	viperInstance.OnConfigChange(func(e fsnotify.Event) {
		err := configuration.Update()
		if err != nil {
			logger.Logger().Error("an error happened when updating the config: " + err.Error())
		}
	})

	return configuration, configuration.Update()
}

func getViperInstance(path string) *viper.Viper {
	viperInstance := viper.New()
	viperInstance.SetConfigFile(path)
	return viperInstance
}

func (configuration *Config) Update() error {
	var cfg ApplicationConfig
	err := configuration.viperInstance.Unmarshal(&cfg)
	if err != nil {
		return err
	}

	configuration.mutex.Lock()
	defer configuration.mutex.Unlock()
	configuration.config = &cfg
	return nil
}

func (configuration *Config) GetConfig() ApplicationConfig {
	configuration.mutex.RLock()
	defer configuration.mutex.RUnlock()
	return *configuration.config
}
