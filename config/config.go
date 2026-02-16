package config

import (
	"log"
	"os"

	model "github.com/deeraj-kumar/exam-audit/domain"
	"github.com/spf13/viper"
)

var Cfg model.Config

func LoadConfig() error {
	workingDir := os.Getenv("WORKING_DIR")
	viper.AddConfigPath(workingDir + "/config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("failed to read config.yaml file , err - %v", err)
		return err
	}

	viper.Unmarshal(&Cfg)
	Cfg.WorkingDir = workingDir
	return nil
}
