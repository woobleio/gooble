package lib

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sqlx.DB

func LoadDB() {
	DB = sqlx.MustOpen("postgres", getDBUrl())
}

func getDBUrl() string {
	var host, dbName, port, username, passwd = viper.GetString("db_host"),
		viper.GetString("db_name"),
		viper.GetString("db_port"),
		viper.GetString("db_username"),
		viper.GetString("db_password")

	var dbUrl = "postgres://"

	if username != "" {
		dbUrl += username
		dbUrl += ":" + passwd + "@"
	}

	dbUrl += host

	if port != "" {
		dbUrl += ":" + port
	}

	dbUrl += "/" + dbName

	dbUrl += "?sslmode=disable"

	return dbUrl
}

type Query struct {
	Q   string
	Opt *Option
}

func (q *Query) String() string {
	q.build()
	return fmt.Sprintf("%s", q.Q)
}

func (q *Query) build() {
	var str string
	if q.Opt.limit > 0 {
		str = strconv.FormatInt(q.Opt.limit, 10)
		q.Q += " LIMIT " + str
	}
	if q.Opt.offset > 0 {
		str = strconv.FormatInt(q.Opt.offset, 10)
		q.Q += " OFFSET " + str
	}
	q.Q += ";"
}

type NullTime struct {
	pq.NullTime
}

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal("")
	}
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Time = *x
	} else {
		v.Valid = false
	}
	return nil
}

// see https://gist.github.com/adharris/4163702
type StringSlice []string

func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

	var (
		unquotedChar  = `[^",\\{}\s(NULL)]`
		unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

		quotedChar  = `[^"\\]|\\"|\\\\`
		quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

		arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

		arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

		valueIndex int
	)

	array := string(asBytes)

	parsed := make([]string, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		parsed = append(parsed, s)
	}

	(*s) = StringSlice(parsed)

	return nil
}

func (s StringSlice) Value() (driver.Value, error) {
	for i, elem := range s {
		s[i] = `"` + strings.Replace(strings.Replace(elem, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(s, ",") + "}", nil
}
