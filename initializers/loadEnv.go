package initializers

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Apikey       string        `mapstructure:"API_KEY"`
	SheetId      string        `mapstructure:"SHEET_ID"`
	Weebhook     string        `mapstructure:"WEBHOOK"`
	DevlopmentId string        `mapstructure:"DEVELOPMENT_ID"`
	ScriptURL    string        `mapstructure:"URL"`
	Datetime     time.Duration `mapstructure:"JWT_EXPIRED_IN"`
	ClientOrigin string        `mapstructure:"CLIENT_ORIGIN"`
	SpreadsheetID string		`mapstructure:"SPREADSHEET_ID"`
	Testsheet string			`mapstructure:"TEST_SHEET"`
	credentials string			`mapstructure:"CREDENTIALS"`
	Solosheet string			`mapstructure:"SOLO_SHEET"`
	Duosheet string				`mapstructure:"DUO_SHEET"`
	Teamsheet string			`mapstructure:"TEAM_SHEET"`

	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`

	
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
