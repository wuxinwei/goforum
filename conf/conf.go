package conf

import (
	"encoding/gob"
	"log"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wuxinwei/goforum/models"
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
	Elastic struct {
		Name     string `mapstructure:"name"`
		Address  string `mapstructure:"address"`
		Port     int    `mapstructure:"port"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
	} `mapstructure:"elastic"`
}

var (
	GlobalConf    ServerConf
	SessionSecret = []byte("secret")
)

func init() {
	// read conf
	viper.AddConfigPath(os.Getenv("HOME"))
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
	// elastic default value
	viper.SetDefault("elastic.name", "lab")
	viper.SetDefault("elastic.address", "127.0.0.1")
	viper.SetDefault("elastic.port", 9200)
	viper.SetDefault("elastic.username", "elastic")
	viper.SetDefault("elastic.password", "elasticpw")

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

	if GlobalConf.Log.Formatter == "text" {
		logrus.SetFormatter(&logrus.TextFormatter{})
	} else {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// register gob type
	gob.Register(&models.User{})
	gob.Register(&models.Post{})
	gob.Register(&models.Comment{})
	gob.Register(&models.Tag{})
}
