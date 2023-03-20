package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DBHost                        string `mapstructure:"POSTGRES_HOST"`
	DBUserName                    string `mapstructure:"POSTGRES_USER"`
	DBUserPassword                string `mapstructure:"POSTGRES_PASSWORD"`
	DBName                        string `mapstructure:"POSTGRES_DB"`
	DBPort                        string `mapstructure:"POSTGRES_PORT"`
	JWT_SECRET                    string `mapstructure:"JWT_SECRET"`
	JWT_EXPIRY_HOUR               string `mapstructure:"JWT_EXPIRY_HOUR"`
	JWT_REFRESH_TOKEN_EXPIRY_HOUR string `mapstructure:"JWT_REFRESH_TOKEN_EXPIRY_HOUR"`
	JWT_REFRESH_TOKEN_SECRET      string `mapstructure:"JWT_REFRESH_TOKEN_SECRET"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func GetDsn(config *Config) string {
	// Used for local
	// dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
	// 	config.DBHost, config.DBUserName, config.DBUserPassword, config.DBName, config.DBPort)

	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", config.DBUserName, config.DBUserPassword, config.DBHost, config.DBName)

	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	log.Fatal("Failed to connect to the Database")
	// }
	// fmt.Println("🚀 Connected Successfully to the Database")
	return dsn
}
