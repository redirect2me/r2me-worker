package main

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigData struct {
	HttpHost    string
	HttpPort    int
	HttpsHost   string
	HttpsPort   int
	Hostname    string
	Action      string
	Endpoint    string
	AcmeEmail   string
	AcmeStaging bool
}

var Config = LoadConfig()

func LoadConfig() *ConfigData {

	pflag.String("http_host", "", "Host to run http server on")
	pflag.Int("http_port", 80, "Port to run http server on")
	pflag.String("https_host", "", "Host to run https server on")
	pflag.Int("https_port", 443, "Port to run https server on")
	pflag.String("hostname", "localhost", "Hostname of this server (for management pages)")
	pflag.String("action", "addwww", "action [lookup|addwww|removewww|api]")
	pflag.String("endpoint", "", "endpoint for api action")
	pflag.String("acme_email", "", "Email for ACME registration")
	pflag.Bool("acme_staging", false, "Use ACME staging environment")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	viper.AutomaticEnv()

	viper.SetConfigName("r2me")
	viper.AddConfigPath("/etc/r2me/")
	viper.AddConfigPath("$HOME/.r2me")
	viper.AddConfigPath(".")
	viper.ReadInConfig()

	config := &ConfigData{
		HttpHost:  viper.GetString("http_host"),
		HttpPort:  viper.GetInt("http_port"),
		HttpsHost: viper.GetString("https_host"),
		HttpsPort: viper.GetInt("https_port"),
		Hostname:  viper.GetString("hostname"),
		Action:    viper.GetString("action"),
		Endpoint:  viper.GetString("endpoint"),
	}

	return config
}
