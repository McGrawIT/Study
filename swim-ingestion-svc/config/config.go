package config

import (
	"encoding/json"
	"github.build.ge.com/aviation-predix-common/vcap-support"
	"github.com/jackmanlabs/errors"
	"log"
	"os"
)

var ActiveConfig *Config

func init() {
	loadDefaultConfig()
	ActiveConfig = DefaultConfig
}

type Config struct {
	PostgreSQL struct {
		Database string
		Host     string
		Password string
		Username string
		Port     string
	}
}

func LoadFromFile(filename string) error {
	f, err := os.Open(filename)
	if os.IsNotExist(err) {
		log.Println("The configuration file was not found.")
		log.Println("A new one is being generated for you.")
		err = writeDefault(filename)
		if err != nil {
			log.Println("An error was encountered while generated your configuration.")
			return errors.Stack(err)
		}

		return errors.New("Please restart with the new configuration file.")
	} else if err != nil {
		return errors.Stack(err)
	}
	defer f.Close()

	ActiveConfig = new(Config)
	err = json.NewDecoder(f).Decode(ActiveConfig)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}

func LoadFromVcap() error {
	// This environment variable is defined in the Predix manifest.
	database := os.Getenv("DB_POSTGRES_NAME")

	var credentials map[string]interface{}

	// Get out the Vcap configuration map.
	// I'd like to clean this up, maybe with static struct definitions?
	vcapServices, err := vcap.LoadServices()
	if err != nil {
		return errors.Stack(err)
	}
	for i := range vcapServices["postgres"] {
		if vcapServices["postgres"][i].Name == database {
			credentials = vcapServices["postgres"][i].Credentials
		}
	}

	if credentials == nil {
		return errors.New("Failure to find database configuration via VCAP.")
	}

	ActiveConfig = &Config{}
	ActiveConfig.PostgreSQL.Database = credentials["database"].(string)
	ActiveConfig.PostgreSQL.Password = credentials["password"].(string)
	ActiveConfig.PostgreSQL.Username = credentials["username"].(string)
	ActiveConfig.PostgreSQL.Host = credentials["host"].(string)
	ActiveConfig.PostgreSQL.Port = credentials["port"].(string)

	return nil
}
