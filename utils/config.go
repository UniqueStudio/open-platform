package utils

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/configor"
)

// Config is a struct for backend config
type Config struct {
	APPName string `default:"Gin App"`

	Server struct {
		Host      string `default:"127.0.0.1"`
		Hostname  string `default:"open.hustunique.com"`
		Port      string `default:"9012"`
		SecretKey string `default:"SecretKey"`
	}

	QcloudSMS struct {
		AppID  string `default:""`
		AppKey string `default:""`
		Sign   string `default:""`
	}

	WeWork struct {
		CropID        string `required:"true"` // CorpID
		AgentID       int    `required:"true"` // Application ID
		AgentSecret   string `required:"true"`
		Secret        string `required:"true"` // Application Secret
		ContactSecret string `required:"true"`
	}

	SMTP struct {
		Sender   string `required:"true"`
		Password string `required:"true"`
		Host     string `required:"true"`
	}
}

// LoadConfiguration is a function to load cfg from file
func LoadConfiguration() Config {
	path, err := os.Getwd()

	switch gin.Mode() {
	case "release":
		path = strings.Replace(path, "test", "", -1) + "/config.deploy.yml"
	case "debug":
		path = strings.Replace(strings.Replace(path, "test", "", -1), "/handler", "", -1) + "/config.yml"
	}

	var config Config
	configFile, err := os.Open(path)
	defer configFile.Close()
	if err != nil {
		log.Printf("[loadAppConfig]: %s\n", err)
	}

	configor.Load(&config, path)

	// Generate secret key
	if count := len([]rune(config.Server.SecretKey)); count <= 32 {
		for i := 1; i <= 32-count; i++ {
			config.Server.SecretKey += "="
		}
	} else {
		config.Server.SecretKey = string([]byte(config.Server.SecretKey)[:32])
	}

	fmt.Printf("%v", config)
	return config
}

// AppConfig is a struct loaded from config file
var AppConfig = LoadConfiguration()
