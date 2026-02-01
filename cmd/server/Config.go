package main

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigData struct {
	AcmeEmail   string
	AcmeStaging bool
	Action      string
	AdminIP     string
	CertDir     string
	Endpoint    string
	Hostname    string
	HttpHost    string
	HttpPort    int
	HttpsHost   string
	HttpsPort   int
	LoadError   error
	LogFormat   string
	LogLevel    string
}

var Config = LoadConfig()

func LoadConfig() *ConfigData {

	pflag.String("acme_email", "", "Email for ACME registration")
	pflag.Bool("acme_staging", true, "Use ACME staging environment")
	pflag.String("action", "addwww", "action [lookup|addwww|removewww|api]")
	pflag.String("admin_ip", "auto", "IP address for direct access to admin pages [ auto | none | specific IP ]")
	pflag.String("cert_dir", os.TempDir(), "Directory to store ACME certs")
	pflag.String("endpoint", "", "endpoint for api action")
	pflag.String("hostname", "localhost", "Hostname of this server (for management pages)")
	pflag.String("http_host", "", "Host to run http server on")
	pflag.Int("http_port", 80, "Port to run http server on")
	pflag.String("https_host", "", "Host to run https server on")
	pflag.Int("https_port", 443, "Port to run https server on")
	pflag.String("log_format", "json", "Log format [ json | text ]")
	pflag.String("log_level", "info", "Log level [ trace | debug | info | warn | error ]")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/redirect2me/")
	viper.AddConfigPath(".")
	LoadError := viper.ReadInConfig()

	config := &ConfigData{
		AcmeEmail:   viper.GetString("acme_email"),
		AcmeStaging: viper.GetBool("acme_staging"),
		Action:      viper.GetString("action"),
		AdminIP:     viper.GetString("admin_ip"),
		CertDir:     viper.GetString("cert_dir"),
		Endpoint:    viper.GetString("endpoint"),
		Hostname:    viper.GetString("hostname"),
		HttpHost:    viper.GetString("http_host"),
		HttpPort:    viper.GetInt("http_port"),
		HttpsHost:   viper.GetString("https_host"),
		HttpsPort:   viper.GetInt("https_port"),
		LoadError:   LoadError,
		LogFormat:   viper.GetString("log_format"),
		LogLevel:    viper.GetString("log_level"),
	}

	return config
}

func (c *ConfigData) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		Logger.Error("ConfigData json.Marshal failed", "error", err)
		b = []byte("{\"success\":false,\"err\":\"json.Marshal failed\"}")
	}
	return string(b)
}

func (c *ConfigData) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Attr{Key: "acme_email", Value: slog.StringValue(c.AcmeEmail)},
		slog.Attr{Key: "acme_staging", Value: slog.BoolValue(c.AcmeStaging)},
		slog.Attr{Key: "action", Value: slog.StringValue(c.Action)},
		slog.Attr{Key: "admin_ip", Value: slog.StringValue(c.AdminIP)},
		slog.Attr{Key: "cert_dir", Value: slog.StringValue(c.CertDir)},
		slog.Attr{Key: "endpoint", Value: slog.StringValue(c.Endpoint)},
		slog.Attr{Key: "hostname", Value: slog.StringValue(c.Hostname)},
		slog.Attr{Key: "http_host", Value: slog.StringValue(c.HttpHost)},
		slog.Attr{Key: "http_port", Value: slog.IntValue(c.HttpPort)},
		slog.Attr{Key: "https_host", Value: slog.StringValue(c.HttpsHost)},
		slog.Attr{Key: "https_port", Value: slog.IntValue(c.HttpsPort)},
		slog.Attr{Key: "log_format", Value: slog.StringValue(c.LogFormat)},
		slog.Attr{Key: "log_level", Value: slog.StringValue(c.LogLevel)},
	)
}
