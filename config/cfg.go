package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Settings struct {
	Application ApplicationConfig `mapstructure:"application"`
	APM         APMConfig         `mapstructure:"apm"`
	Database    DatabaseConfig    `mapstructure:"database"`
	Email       EmailConfig       `mapstructure:"email"`
	Tencent     TencentConfig     `mapstructure:"tencent"`
}

type ApplicationConfig struct {
	Name         string `mapstructure:"name"`
	Host         string `mapstructure:"host"`
	HttpPort     string `mapstructure:"http_port"`
	GrpcPort     string `mapstructure:"grpc_port"`
	GrpcCertFile string `mapstructure:"grpc_cert"`
	GrpcKeyFile  string `mapstructure:"grpc_key"`
	Hostname     string `mapstructure:"hostname"`
	Mode         string `mapstructure:"mode"`
	SecretKey    string `mapstructure:"secret_key"`
	SessionKey   string `mapstructure:"session_key"`
}

type APMConfig struct {
	ReporterBackground string `mapstructure:"reporter_backend"`
}

type DatabaseConfig struct {
	Postgres struct {
		UserDSN string `mapstructure:"user_dsn"`
		OpenDSN string `mapstructure:"open_dsn"`
	} `mapstructure:"postgres"`
	Redis struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	}
}

type TencentConfig struct {
	SMS SMSConfig `mapstructure:"sms"`
}

type SMSConfig struct {
	SDKAppID             string `mapstructure:"sdk_app_id"`
	SecretID             string `mapstructure:"secret_id"`
	SecretKey            string `mapstructure:"secret_key"`
	DefaultVirtualSignId uint   `mapstructure:"default_virtual_sign_id"`
}

type EmailConfig struct {
	SMTP          SMTPConfig            `mapstructure:"smtp"`
	EmailTemplate []EmailTemplateConfig `mapstructure:"email_template"`
}

type SMTPConfig struct {
	Sender   string `mapstructure:"sender"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Name     string `mapstructure:"name"`
	Port     int    `mapstructure:"port"`
}

type EmailTemplateConfig struct {
	TemplateID string `mapstructure:"template_id"`
	Filename   string `mapstructure:"filename"`
}

var Config *Settings

func SetupConfig(filepath string) error {
	Config = new(Settings)
	viper.SetConfigFile(filepath)
	viper.SetDefault("application.host", "0.0.0.0")
	viper.SetDefault("application.name", "open-platform")
	viper.SetDefault("application.port", "5000")
	viper.SetDefault("application.hostname", "https://open.hustunique.com")

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	err = viper.Unmarshal(Config)
	if err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(Config); err != nil {
		return err
	}

	return nil
}
