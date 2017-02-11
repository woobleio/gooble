package lib

import (
	"database/sql"
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

// DB driver
var DB *sqlx.DB

// LoadDB initializes the DB
func LoadDB() {
	DB = sqlx.MustOpen("postgres", getDBUrl())
}

func getDBUrl() string {
	var host, dbName, port, username, passwd = viper.GetString("db_host"),
		viper.GetString("db_name"),
		viper.GetString("db_port"),
		viper.GetString("db_username"),
		viper.GetString("db_password")

	var dbURL = "postgres://"

	if username != "" {
		dbURL += username
		dbURL += ":" + passwd + "@"
	}

	dbURL += host

	if port != "" {
		dbURL += ":" + port
	}

	dbURL += "/" + dbName

	dbURL += "?sslmode=disable"

	return dbURL
}

// Query is all all query options
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
	if q.Opt.Limit > 0 {
		str = strconv.FormatInt(q.Opt.Limit, 10)
		q.Q += " LIMIT " + str
	}
	if q.Opt.Offset > 0 {
		str = strconv.FormatInt(q.Opt.Offset, 10)
		q.Q += " OFFSET " + str
	}
	q.Q += ";"
}

// NullTime is psql null for time.Time
type NullTime struct {
	pq.NullTime
}

// InitNullTime returns a NullTime with Time "date"
func InitNullTime(date time.Time) *NullTime {
	return &NullTime{
		pq.NullTime{
			Time:  date,
			Valid: true,
		},
	}
}

// MarshalJSON marshals custom NullTime
func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	}
	return json.Marshal("")
}

// UnmarshalJSON unmarshals custom NullTime
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

// NullString is psql null for string
type NullString struct {
	sql.NullString
}

// InitNullString returns a NullString with String "str"
func InitNullString(str string) *NullString {
	return &NullString{
		sql.NullString{
			String: str,
			Valid:  true,
		},
	}
}

// MarshalJSON marshals custom NullString
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal("")
}

// NullInt64 is psql null for string
type NullInt64 struct {
	sql.NullInt64
}

// InitNullInt64 returns a NullInt64 with int64 "val"
func InitNullInt64(val int64) *NullInt64 {
	return &NullInt64{
		sql.NullInt64{
			Int64: val,
			Valid: true,
		},
	}
}

// MarshalJSON marshals custom NullString
func (ns NullInt64) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.Int64)
	}
	return json.Marshal("")
}

// UnmarshalJSON unmarshals custom NullString
func (ns *NullInt64) UnmarshalJSON(data []byte) error {
	var x int64
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}

	// TODO control if it's a number
	ns.Valid = true
	ns.Int64 = x

	return nil
}

// ID is a custom type for hashed IDs
type ID struct {
	ValueEncoded string
	ValueDecoded int64
}

// MarshalJSON marshals custom ID
func (id ID) MarshalJSON() ([]byte, error) {
	return json.Marshal(id.ValueEncoded)
}

// UnmarshalJSON unmarshals custom NullString
func (id *ID) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		id.ValueEncoded = *x
		decode, err := DecodeHash(id.ValueEncoded)
		if err != nil {
			return err
		}
		id.ValueDecoded = decode
	} else {
		id.ValueEncoded = ""
		id.ValueDecoded = 0
	}
	return nil
}

// Scan implements the Scanner interface
func (id *ID) Scan(value interface{}) error {
	id.ValueDecoded = value.(int64)
	encode, err := HashID(id.ValueDecoded)
	if err != nil {
		return err
	}
	id.ValueEncoded = encode
	return nil
}

// Value implements the driver Valuer interface
func (id ID) Value() (driver.Value, error) {
	return id.ValueDecoded, nil
}

// StringSlice see https://gist.github.com/adharris/4163702
type StringSlice []string

// Scan scans StringSlice
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
		s = strings.Trim(s, ",")
		parsed = append(parsed, s)
	}

	(*s) = StringSlice(parsed)

	return nil
}

// Value returns StringSlice as psql value
func (s StringSlice) Value() (driver.Value, error) {
	for i, elem := range s {
		s[i] = `"` + strings.Replace(strings.Replace(elem, `\`, `\\\`, -1), `"`, `\"`, -1) + `"`
	}
	return "{" + strings.Join(s, ",") + "}", nil
}
