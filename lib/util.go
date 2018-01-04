package lib

import (
	"encoding/hex"
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"math"
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
	hasher, _ := initHasher()
	return hasher.EncodeInt64([]int64{id})
}

// DecodeHash returns the decoded id
func DecodeHash(hash string) (int64, error) {
	hasher, _ := initHasher()
	id, err := hasher.DecodeInt64WithError(hash)
	if err != nil {
		return 0, err
	}
	return id[0], nil
}

type circle struct {
	X, Y, R float64
}

type pattern struct {
	A      float64
	B      float64
	K      float64
	T      float64
	Seq    int
	RGB    []uint8
	RGBRat []int
}

// GenImage generates an almost unique cool image
func GenImage(id uint64) *image.RGBA {
	var w, h int = 240, 240
	var hw, hh float64 = float64(w / 2), float64(h / 2)
	cr := &circle{hw - math.Sin(0), hh - math.Cos(0), 80}

	m := image.NewRGBA(image.Rect(0, 0, w, h))
	pat := newPattern(id)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			var dx, dy float64 = cr.X - float64(x), cr.Y - float64(y)
			c := color.RGBA{0, 0, 0, 0}
			ar := math.Sqrt(dx*dx+dy*dy) / cr.R
			t := float64(x) * pat.T
			x2 := pat.B*t + pat.A*pat.K*math.Sin(0.5*t)
			y2 := (-pat.B*t + pat.B*pat.K*math.Sin(0.5*t))

			if x2+240 >= float64(x) && ar < 1.5 {
				var curSeq = h
				for curSeq > 0 {
					curSeq -= pat.Seq
					if y2+float64(curSeq) < float64(y) {
						r := ((pat.RGB[0] + uint8(curSeq*pat.RGBRat[0])) % 155) + 100
						g := ((pat.RGB[1] + uint8(curSeq*pat.RGBRat[1])) % 155) + 50
						b := ((pat.RGB[2] + uint8(curSeq*pat.RGBRat[2])) % 155) + 90
						c = color.RGBA{r, g, b, 255}
						break
					}
				}
			}

			m.Set(x, y, c)
		}
	}

	return m
}

func initHasher() (*hashids.HashID, error) {
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

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func newPattern(id uint64) *pattern {
	u := hash(fmt.Sprintf("%d", id))
	k := (u % 20) + 1
	t := u%233 + 50
	a := math.Sin(0)
	b := math.Cos(0)
	c := uint8((u * t) % 255)
	return &pattern{
		A:      a,
		B:      b,
		K:      float64(k),
		T:      float64(float64(t) / 1000),
		Seq:    40,
		RGB:    []uint8{c, c, c},
		RGBRat: []int{5, 3, 4},
	}
}
