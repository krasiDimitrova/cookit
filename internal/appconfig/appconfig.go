package appconfig

import (
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
)

// appConfig provides db and project config values
type appConfig struct {
	db_name     string
	db_password string
	db_username string
	db_host     string
	project_dir string
}

// AppConfig interface provide methods for obtaining config values
type AppConfig interface {
	// GetDBConfig function returns db connection information
	GetDBConfig() (username string, password string, databaseName string, databaseHost string)

	// GetProjectDir function returns the project directory path
	GetProjectDir() (projectDir string)
}

var config appConfig
var so sync.Once

func Get() AppConfig {
	so.Do(func() {
		config = appConfig{}
		config.loadConfiguration()
	})

	return &config
}

func (config *appConfig) GetDBConfig() (username string, password string, databaseName string, databaseHost string) {
	username = config.db_username
	password = config.db_password
	databaseName = config.db_name
	databaseHost = config.db_host
	return
}

func (config *appConfig) GetProjectDir() (projectDir string) {
	return config.project_dir
}

func (config *appConfig) loadConfiguration() {
	config.project_dir, _ = os.Getwd()

	viper.SetConfigName("app")
	// Set the path to look for the configurations file
	viper.AddConfigPath(config.project_dir + "/configs")
	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	config.db_name = viper.GetString("MYSQL_DATABASE")
	config.db_host = viper.GetString("MYSQL_SERVICE_HOST")
	config.db_username = viper.GetString("MYSQL_USERNAME")
	config.db_password = viper.GetString("MYSQL_PASSWORD")

	return
}
