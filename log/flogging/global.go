/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package flogging

import (
	"io"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/grpclog"
)

const (
	defaultFormat = "json"
	defaultLevel  = zapcore.DebugLevel

	configKey = "CONFIG_FILE"
)

var (
	Global *Logging
	config *viper.Viper
)

func init() {
	logging, err := New(Config{})
	if err != nil {
		panic(err)
	}

	Global = logging
	grpcLogger := Global.ZapLogger("grpc")
	grpclog.SetLoggerV2(NewGRPCLogger(grpcLogger))

	loadConfig()
}

// loadConfig loads config file with path get from CONFIG_FILE environment variable
//	by default, config file will have path as PWD/config/config.yaml
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
			logLevel := defaultLevel.String()
			logFormat := defaultFormat
			cfgMap := config.GetStringMap("Log")
			if val, ok := cfgMap[strings.ToLower("Level")]; ok {
				logLevel = val.(string)
			}
			if val, ok := cfgMap[strings.ToLower("Format")]; ok {
				logFormat = val.(string)
			}
			log.Printf("Config log level to %s, format as %s", logLevel, logFormat)

			Global.Apply(Config{
				Format:  logFormat,
				LogSpec: logLevel,
			})
		}

	}
}

// Init initializes logging with the provided config.
func Init(config Config) {
	err := Global.Apply(config)
	if err != nil {
		panic(err)
	}
}

// Reset sets logging to the defaults defined in this package.
//
// Used in tests and in the package init
func Reset() {
	Global.Apply(Config{})
}

// LoggerLevel gets the current logging level for the logger with the
// provided name.
func LoggerLevel(loggerName string) string {
	return Global.Level(loggerName).String()
}

// MustGetLogger creates a logger with the specified name. If an invalid name
// is provided, the operation will panic.
func MustGetLogger(loggerName string) *FabricLogger {
	return Global.Logger(loggerName)
}

// ActivateSpec is used to activate a logging specification.
func ActivateSpec(spec string) {
	err := Global.ActivateSpec(spec)
	if err != nil {
		panic(err)
	}
}

// DefaultLevel returns the default log level.
func DefaultLevel() string {
	return defaultLevel.String()
}

// SetWriter calls SetWriter returning the previous value
// of the writer.
func SetWriter(w io.Writer) io.Writer {
	return Global.SetWriter(w)
}

// SetObserver calls SetObserver returning the previous value
// of the observer.
func SetObserver(observer Observer) Observer {
	return Global.SetObserver(observer)
}
