package lib

import (
	"math/rand"

	"github.com/spf13/viper"
)

func GenKey() string {
	keyRange := viper.GetString("key_range")
	b := make([]byte, 15)
	for i := range b {
		b[i] = keyRange[rand.Intn(len(keyRange))]
	}
	return string(b)
}

func GetTokenLifetime() int {
	return viper.GetInt("token_lifetime")
}
