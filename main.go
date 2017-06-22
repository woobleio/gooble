package main

import (
	"math/rand"
	"os"
	"time"

	"wooble/lib"
	"wooble/router"

	"github.com/spf13/viper"
	stripe "github.com/stripe/stripe-go"
)

func init() {
	if os.Getenv("GOENV") == "prod" {
		stripe.Key = os.Getenv("STRIPE_KEY")
	} else {
		viper.SetConfigName(os.Getenv("GOENV"))
		viper.SetConfigType("yaml")
		viper.AddConfigPath(os.Getenv("CONFPATH"))

		if err := viper.ReadInConfig(); err != nil {
			panic(err)
		}
		stripe.Key = viper.GetString("stripe_key")
	}

	lib.LoadDB()
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	router.Load()
}
