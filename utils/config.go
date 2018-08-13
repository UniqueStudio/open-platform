package utils

import (
	"fmt"
	"log"

	"github.com/denverdino/aliyungo/acm"
	"github.com/gin-gonic/gin"
	yaml "gopkg.in/yaml.v2"
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

	const AccessKeyID = "LTAIZcjgmcBs7Y95"
	const AccessKeySecret = "cS6XBX5uIeLDcfRtQD6pgj2gK9Dq0e"

	client, err := acm.NewClient(func(c *acm.Client) {
		c.AccessKey = AccessKeyID
		c.SecretKey = AccessKeySecret
		c.EndPoint = "acm.aliyun.com"
		c.NameSpace = "0fb30cff-c36f-495a-a76e-9ab434994739"
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", gin.Mode())
	ret, err := client.GetConfig("studio.open-platform."+gin.Mode(), "STUDIO")
	if err != nil {
		log.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal([]byte(ret), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

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
