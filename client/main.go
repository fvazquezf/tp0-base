package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
)

// InitConfig Function that uses viper library to parse configuration parameters.
// Viper is configured to read variables from both environment variables and the
// config file ./config.yaml. Environment variables takes precedence over parameters
// defined in the configuration file. If some of the variables cannot be parsed,
// an error is returned
func InitConfig() (*viper.Viper, error) {
	v := viper.New()

	// Configure viper to read env variables with the CLI_ prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("cli")
	// Use a replacer to replace env variables underscores with points. This let us
	// use nested configurations in the config file and at the same time define
	// env variables for the nested configurations
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Add env variables supported
	v.BindEnv("agency")
	v.BindEnv("server", "address")
	v.BindEnv("log", "level")

	// Add clients user vaiables
	v.BindEnv("name")
	v.BindEnv("lastname")
	v.BindEnv("birthdate")
	v.BindEnv("id")
	v.BindEnv("number")

	// Try to read configuration from config file. If config file
	// does not exists then ReadInConfig will fail but configuration
	// can be loaded from the environment variables so we shouldn't
	// return an error in that case
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file. Using env variables instead")
	}

	return v, nil
}

// InitLogger Receives the log level to be set in logrus as a string. This method
// parses the string and set the level to the logger. If the level string is not
// valid an error is returned
func InitLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

    customFormatter := &logrus.TextFormatter{
      TimestampFormat: "2006-01-02 15:04:05",
      FullTimestamp: false,
    }
    logrus.SetFormatter(customFormatter)
	logrus.SetLevel(level)
	return nil
}

// PrintConfig Print all the configuration parameters of the program.
// For debugging purposes only
func PrintConfig(v *viper.Viper) {
	logrus.Infof("action: config | result: success | client_agency: %s | server_address: %s | log_level: %s",
	    v.GetString("agency"),
	    v.GetString("server.address"),
	    v.GetString("log.level"),
    )
}

func main() {
	v, err := InitConfig()
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := InitLogger(v.GetString("log.level")); err != nil {
		log.Fatalf("%s", err)
	}

	// Print program config with debugging purposes
	PrintConfig(v)

	clientConfig := common.ClientConfig{
		ServerAddress: v.GetString("server.address"),
		Agency:        v.GetString("agency"),
		LoopLapse:     v.GetDuration("loop.lapse"),
		LoopPeriod:    v.GetDuration("loop.period"),
	}

	bet := &common.Bet{
			Name:      v.GetString("name"),
			LastName:  v.GetString("lastname"),
			BirthDate: v.GetString("birthdate"),
			ID:        v.GetString("id"),
			Number:    uint16(v.GetUint("number")),
			AgencyId:  uint16(v.GetUint("agency")),
		}

	client := common.NewClient(clientConfig, bet)
	logrus.Infof("action: NewClient | result: success")
	client.SendBetAndValidate()
}
