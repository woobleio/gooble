package lib

import (
	"encoding/hex"
	"math/rand"

	"golang.org/x/crypto/scrypt"

	hashids "github.com/speps/go-hashids"
	"github.com/spf13/viper"
)

// GenKey generates random key
func GenKey() string {
	keyRange := viper.GetString("keygen_range")
	b := make([]byte, 15)
	for i := range b {
		b[i] = keyRange[rand.Intn(len(keyRange))]
	}
	return string(b)
}

// Encrypt encrypt the string
func Encrypt(toEncrypt string, salt []byte) (string, error) {
	cp, err := scrypt.Key([]byte(toEncrypt), []byte(salt), 16384, 8, 1, 32)
	return hex.EncodeToString(cp), err
}

// GetTokenLifetime returns token lifetime
func GetTokenLifetime() int {
	return viper.GetInt("token_lifetime")
}

// GetOrigins returns origins for cors
func GetOrigins() []string {
	return viper.GetStringSlice("allow_origin")
}

// GetPkgURL returns packages URL
func GetPkgURL() string {
	return viper.GetString("pkg_url")
}

// HashID hashes uint64 id and returns a unique string
func HashID(id int64) (string, error) {
	hasher := initHasher()
	return hasher.EncodeInt64([]int64{id})
}

// DecodeHash returns the decoded id
func DecodeHash(hash string) (int64, error) {
	hasher := initHasher()
	id, err := hasher.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}
	return id[0], nil
}

func initHasher() *hashids.HashID {
	hashConf := hashids.NewData()
	hashConf.Salt = viper.GetString("salt_for_id")
	hashConf.MinLength = 8
	return hashids.NewWithData(hashConf)
}
