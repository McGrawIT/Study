package config

import (
	"encoding/json"
	"github.com/jackmanlabs/errors"
	"os"
)

var DefaultConfig *Config

func loadDefaultConfig(){
	DefaultConfig = &Config{}
	DefaultConfig.PostgreSQL.Database = "swim-ingestion-svc"
	DefaultConfig.PostgreSQL.Host = "localhost"
	DefaultConfig.PostgreSQL.Password = "swim-ingestion-svc"
	DefaultConfig.PostgreSQL.Username = "swim-ingestion-svc"
}


func writeDefault(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return errors.Stack(err)
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(DefaultConfig)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
