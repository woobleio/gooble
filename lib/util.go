package lib

import (
	"math/rand"

	"github.com/spf13/viper"
)

// GenKey generates random key
func GenKey() string {
	keyRange := viper.GetString("key_range")
	b := make([]byte, 15)
	for i := range b {
		b[i] = keyRange[rand.Intn(len(keyRange))]
	}
	return string(b)
}

// GetTokenLifetime returns token lifetime
func GetTokenLifetime() int {
	return viper.GetInt("token_lifetime")
}

// GetOrigins returns origins for cors
func GetOrigins() []string {
	return viper.GetStringSlice("allow_origin")
}
