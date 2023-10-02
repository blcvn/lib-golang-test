package config

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init() {
	var err error

	configFile := os.Getenv("CONFIG_FILE")
	if configFile == "" {
		configFile = "api.yaml"
	}
	config = viper.New()
	config.SetConfigType("yaml")

	configDir := path.Dir(configFile)
	fileName := path.Base(configFile)

	config.SetConfigName(fileName)
	config.AddConfigPath(configDir)
	config.AddConfigPath("./config/")
	config.AddConfigPath("../config/")
	config.AddConfigPath("../../config/")

	err = config.ReadInConfig()
	if err != nil {
		log.Fatal("error on parsing configuration file")
	}
}

func relativePath(basedir string, path *string) {
	p := *path
	if len(p) > 0 && p[0] != '/' {
		*path = filepath.Join(basedir, p)
	}
}

func GetConfig() *viper.Viper {
	return config
}

func GatewayPort() int {
	return 8001
}

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
