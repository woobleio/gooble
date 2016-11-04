package main

import (
	"os"

	"wooble/lib"
	"wooble/router"

	"github.com/spf13/viper"
)

func init() {
	viper.SetConfigName(os.Getenv("GOENV"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("CONFPATH"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	lib.LoadDB()
}

func main() {
	router.Load()
}
