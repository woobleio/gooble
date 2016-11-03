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

	errViper := viper.ReadInConfig()
	if errViper != nil {
		panic(errViper)
	}

	lib.LoadDB()
}

func main() {
	router.Load()
}
