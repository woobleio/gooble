package lib

import (
	"encoding/hex"
	"math/rand"
	"os"
	"strconv"

	"golang.org/x/crypto/scrypt"

	hashids "github.com/speps/go-hashids"
	"github.com/spf13/viper"
)

// GenKey generates random key
func GenKey() string {
	var keyRange string
	if isProd() {
		keyRange = os.Getenv("KEYGEN_RANGE")
	} else {
		keyRange = viper.GetString("keygen_range")
	}
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
	if isProd() {
		time, err := strconv.Atoi(os.Getenv("TOKEN_LIFETIME"))
		if err != nil {
			panic(err)
		}
		return time
	}

	return viper.GetInt("token_lifetime")
}

// GetTokenKey returns the token's salt key
func GetTokenKey() string {
	if isProd() {
		return os.Getenv("TOKEN_KEY")
	}

	return viper.GetString("token_key")
}

// GetEncKey returns a key for encryptions
func GetEncKey() string {
	if isProd() {
		return os.Getenv("ENC_KEY")
	}
	return viper.GetString("enc_key")
}

// GetCloudRepo returns the cloud repository name
func GetCloudRepo() string {
	if isProd() {
		return os.Getenv("CLOUD_REPO")
	}

	return viper.GetString("cloud_repo")
}

// GetOrigins returns origins for cors
func GetOrigins() []string {
	if isProd() {
		return []string{os.Getenv("ALLOW_ORIGIN")}
	}

	return viper.GetStringSlice("allow_origin")
}

// GetPkgURL returns packages URL
func GetPkgURL() string {
	if isProd() {
		return os.Getenv("PKG_URL")
	}

	return viper.GetString("pkg_url")
}

// GetPkgRepo returns the package repository
func GetPkgRepo() string {
	if isProd() {
		return os.Getenv("PKG_REPO")
	}

	return viper.GetString("pkg_repo")
}

// GetEmailHost returns email host with port
func GetEmailHost() string {
	if isProd() {
		return os.Getenv("EMAIL_HOST")
	}
	return viper.GetString("email_host")
}

// GetEmailPasswd returns email password
func GetEmailPasswd() string {
	if isProd() {
		return os.Getenv("EMAIL_PASSWD")
	}
	return viper.GetString("email_passwd")
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
	if isProd() {
		hashConf.Salt = os.Getenv("SALT_FOR_ID")
	} else {
		hashConf.Salt = viper.GetString("salt_for_id")
	}
	hashConf.MinLength = 8
	return hashids.NewWithData(hashConf)
}

func isProd() bool {
	return os.Getenv("GOENV") == "prod"
}
