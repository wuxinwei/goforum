package conf

import (
	"log"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// LogLevel is a map to provide a mapping relation
// between log level section in conf file
// and log level which is defined at logrus
var LogLevel = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
}

// ServerConf is a configuration file
type ServerConf struct {
	Log struct {
		Formatter string `mapstructure:"formatter"`
		Level     string `mapstructure:"level"`
	} `mapstructure:"log"`
	DB struct {
		Name     string `mapstructure:"name"`
		Address  string `mapstructure:"address"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Debug    bool   `mapstructure:"debug"`
	} `mapstructure:"db"`
	Cache struct {
		Name     string `mapstructure:"name"`
		Address  string `mapstructure:"address"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"cache"`
}

var GlobalConf ServerConf

func init() {
	// read conf
	viper.AddConfigPath("./conf")
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	// log default value
	viper.SetDefault("log.formatter", "text")
	viper.SetDefault("log.level", "error")
	// db default value
	viper.SetDefault("db.name", "lab")
	viper.SetDefault("db.address", "127.0.0.1")
	viper.SetDefault("db.port", 3306)
	viper.SetDefault("db.username", "test")
	viper.SetDefault("db.password", "test")
	viper.SetDefault("db.debug", false)
	// cache default value
	viper.SetDefault("cache.name", "lab")
	viper.SetDefault("cache.address", "127.0.0.1")
	viper.SetDefault("cache.port", 6379)
	viper.SetDefault("cache.username", "test")
	viper.SetDefault("cache.password", "test")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Read configuration file failed, Err: %s", err)
	} else if err := viper.Unmarshal(&GlobalConf); err != nil {
		log.Fatalf("Read configuration file failed, Err: %s", err)
	}

	// log
	if logLevel, ok := LogLevel[GlobalConf.Log.Level]; !ok {
		log.Printf("invalid logrus level, set as default log level: error")
		logrus.SetLevel(logrus.ErrorLevel)
	} else {
		logrus.SetLevel(logLevel)
	}

	// TODO: 是否需要做检查
	if GlobalConf.Log.Formatter == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
}
