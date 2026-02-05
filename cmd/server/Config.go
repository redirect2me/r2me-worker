package main

import (
	"encoding/json"
	"log/slog"
	"os"

	haikunator "github.com/atrox/haikunatorgo/v2"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type ConfigData struct {
	AcmeEmail   string
	AcmeStaging bool
	Action      string
	AdminHost   string
	AdminIP     string
	CertDir     string
	Endpoint    string
	HttpAddr    string
	HttpPort    int
	HttpsAddr   string
	HttpsPort   int
	LoadError   error
	LogFormat   string
	LogLevel    string
	LogSource   bool
	NodeID      string
}

var Config = LoadConfig()

func LoadConfig() *ConfigData {

	pflag.String("acme_email", "", "Email for ACME registration")
	pflag.Bool("acme_staging", true, "Use ACME staging environment")
	pflag.String("action", "addwww", "action [lookup|addwww|removewww|api]")
	pflag.String("admin_host", "localhost", "Hostname for management pages (none to disable)")
	pflag.String("admin_ip", "auto", "IP address for direct access to admin pages [ auto | none | specific IP ]")
	pflag.String("cert_dir", os.TempDir(), "Directory to store ACME certs")
	pflag.String("endpoint", "", "endpoint for api action")
	pflag.String("http_addr", "", "Network interface for the http server (empty = all interfaces)")
	pflag.Int("http_port", 80, "Port for the http server")
	pflag.String("https_addr", "", "Network interface for the https server (empty = all interfaces)")
	pflag.Int("https_port", 443, "Port for the https server")
	pflag.String("log_format", "json", "Log format [ json | text ]")
	pflag.String("log_level", "info", "Log level [ trace | debug | info | warn | error ]")
	pflag.Bool("log_source", false, "Include source file and line number in logs")
	pflag.String("node_id", "", "Node ID for this server instance")
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
		AdminHost:   viper.GetString("admin_host"),
		AdminIP:     viper.GetString("admin_ip"),
		CertDir:     viper.GetString("cert_dir"),
		Endpoint:    viper.GetString("endpoint"),
		HttpAddr:    viper.GetString("http_addr"),
		HttpPort:    viper.GetInt("http_port"),
		HttpsAddr:   viper.GetString("https_addr"),
		HttpsPort:   viper.GetInt("https_port"),
		LoadError:   LoadError,
		LogFormat:   viper.GetString("log_format"),
		LogLevel:    viper.GetString("log_level"),
		LogSource:   viper.GetBool("log_source"),
		NodeID:      viper.GetString("node_id"),
	}

	if config.NodeID == "" {
		config.NodeID = haikunator.New().Haikunate()
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
		slog.Attr{Key: "admin_host", Value: slog.StringValue(c.AdminHost)},
		slog.Attr{Key: "admin_ip", Value: slog.StringValue(c.AdminIP)},
		slog.Attr{Key: "cert_dir", Value: slog.StringValue(c.CertDir)},
		slog.Attr{Key: "endpoint", Value: slog.StringValue(c.Endpoint)},
		slog.Attr{Key: "http_addr", Value: slog.StringValue(c.HttpAddr)},
		slog.Attr{Key: "http_port", Value: slog.IntValue(c.HttpPort)},
		slog.Attr{Key: "https_addr", Value: slog.StringValue(c.HttpsAddr)},
		slog.Attr{Key: "https_port", Value: slog.IntValue(c.HttpsPort)},
		slog.Attr{Key: "log_format", Value: slog.StringValue(c.LogFormat)},
		slog.Attr{Key: "log_level", Value: slog.StringValue(c.LogLevel)},
		slog.Attr{Key: "log_source", Value: slog.BoolValue(c.LogSource)},
		slog.Attr{Key: "node_id", Value: slog.StringValue(c.NodeID)},
	)
}
