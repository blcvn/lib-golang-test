package smallstep

import (
	"os"
	"path"
	"strings"

	"github.com/Hnampk/prometheuslog/flogging"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
)

const (
	defaultFormat = "json"
	defaultLevel  = zapcore.DebugLevel

	configKey = "CONFIG_FILE"
)

var (
	config *viper.Viper

	CAURL, CertValidDuration string

	configLogger = flogging.MustGetLogger("libs.ca.smallstep.config")
)

func init() {
	loadConfig()
}

func loadConfig() {
	if config == nil {
		file := "config.yaml"
		if value := os.Getenv(configKey); value != "" {
			file = value
		}
		fileName := path.Base(file)

		config = viper.New()
		config.SetConfigType("yaml")
		config.SetConfigName(fileName)
		folder := path.Dir(file)
		config.AddConfigPath(folder)
		err := config.ReadInConfig()
		if err == nil {
			cfgMap := config.GetStringMap("CA")
			if val, ok := cfgMap[strings.ToLower("URL")]; ok {
				CAURL = val.(string)
			}
			if val, ok := cfgMap[strings.ToLower("CertValidDuration")]; ok {
				CertValidDuration = val.(string)
			}
		}

		configLogger.Infof("CA URL: %s", CAURL)
		configLogger.Infof("CA new certificate valid duration: %s", CertValidDuration)
	}
}
