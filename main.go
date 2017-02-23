package main

import (
	"os"

	"wooble/lib"
	"wooble/router"

	"github.com/spf13/viper"
	stripe "github.com/stripe/stripe-go"
)

func init() {
	viper.SetConfigName(os.Getenv("GOENV"))
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("CONFPATH"))

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	lib.LoadDB()

	stripe.Key = viper.GetString("stripe_key")
}

func main() {
	router.Load()
}
