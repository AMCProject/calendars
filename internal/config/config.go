package config

import (
	"github.com/joho/godotenv"
	"os"
)

var Config Configuration

// Configuration settings of the service.
type Configuration struct {
	// Host --> Default 0.0.0.0
	Host string `mapstructure:"HOST" json:"host" default:"0.0.0.0"`
	// Port --> Default 49100
	Port string `mapstructure:"PORT" json:"port" default:"3300"`
	// DBName --> Name of the database. Default "amc.db"
	DBName string `mapstructure:"DB_NAME" json:"DBName" default:"amc.db"`
	// UsersURL --> URL of the users microservice
	UsersURL string `mapstructure:"USERS_URL" json:"UsersURL" default:"0.0.0.0:3100"`
	// MealsURL --> URL of the meals microservice
	MealsURL string `mapstructure:"MEALS_URL" json:"MealsURL" default:"0.0.0.0:3200"`
}

func LoadConfiguration() error {
	err := godotenv.Load("./internal/config/.env")
	if err != nil {
		return err
	}
	Config.Host = os.Getenv("HOST")
	Config.Port = os.Getenv("PORT")
	Config.DBName = os.Getenv("DB_NAME")
	Config.UsersURL = os.Getenv("USERS_URL")
	Config.MealsURL = os.Getenv("MEALS_URL")

	return nil
}
