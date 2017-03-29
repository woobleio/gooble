package lib

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
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
func (ni NullInt64) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int64)
	}
	return json.Marshal("")
}

// UnmarshalJSON unmarshals custom NullInt64
func (ni *NullInt64) UnmarshalJSON(data []byte) error {
	s := string(data)
	if v, err := strconv.Atoi(s); err != nil {
		ni.Int64 = int64(v)
		ni.Valid = true
		return nil
	}
	ni.Int64 = -1
	if s == "null" {
		ni.Valid = true
		return nil
	}
	ni.Valid = false
	return errors.New("Invalid NullINt64: " + s)
}

// ID is a custom type for hashed IDs
type ID struct {
	ValueEncoded string
	ValueDecoded uint64
}

// InitID returns an ID
func InitID(data interface{}) ID {
	switch data.(type) {
	case string:
		encV := data.(string)
		decV, _ := DecodeHash(encV)
		return ID{
			ValueEncoded: encV,
			ValueDecoded: uint64(decV),
		}
	case uint64:
		decV := data.(uint64)
		encV, _ := HashID(int64(decV))
		return ID{
			ValueEncoded: encV,
			ValueDecoded: decV,
		}
	}
	return ID{
		ValueEncoded: "",
		ValueDecoded: 0,
	}
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
		id.ValueDecoded = uint64(decode)
	} else {
		id.ValueEncoded = ""
		id.ValueDecoded = 0
	}
	return nil
}

// Scan implements the Scanner interface
func (id *ID) Scan(value interface{}) error {
	encode, err := HashID(value.(int64))
	if err != nil {
		return err
	}
	id.ValueDecoded = uint64(value.(int64))
	id.ValueEncoded = encode
	return nil
}

// Value implements the driver Valuer interface
func (id ID) Value() (driver.Value, error) {
	return int64(id.ValueDecoded), nil
}

// Exprs for psql arrays
var (
	unquotedChar  = `[^",\\{}\s(NULL)]`
	unquotedValue = fmt.Sprintf("(%s)+", unquotedChar)

	quotedChar  = `[^"\\]|\\"|\\\\`
	quotedValue = fmt.Sprintf("\"(%s)*\"", quotedChar)

	arrayValue = fmt.Sprintf("(?P<value>(%s|%s))", unquotedValue, quotedValue)

	arrayExp = regexp.MustCompile(fmt.Sprintf("((%s)(,)?)", arrayValue))

	valueIndex int
)

// UintSlice if a uint slice for sql driver
type UintSlice []uint64

// Value returns UintSlice as psql value
func (us UintSlice) Value() (driver.Value, error) {
	uints := make([]string, len(us))
	for i, v := range us {
		uints[i] = strconv.FormatUint(v, 10)
	}
	return "{" + strings.Join(uints, ",") + "}", nil
}

// Scan scans UintSlice
func (us *UintSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

	array := string(asBytes)

	parsed := make([]uint64, 0)
	matches := arrayExp.FindAllStringSubmatch(array, -1)
	for _, match := range matches {
		s := match[valueIndex]
		// the string _might_ be wrapped in quotes, so trim them:
		s = strings.Trim(s, "\"")
		s = strings.Trim(s, ",")

		val, _ := strconv.ParseUint(s, 10, 64)
		parsed = append(parsed, val)
	}

	(*us) = UintSlice(parsed)

	return nil
}

// StringSlice see https://gist.github.com/adharris/4163702
type StringSlice []string

// Scan scans StringSlice
func (s *StringSlice) Scan(src interface{}) error {
	asBytes, ok := src.([]byte)
	if !ok {
		return errors.New("Scan source was not []bytes")
	}

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

// SQLPatch is a struct for parsed resource to be patched
type SQLPatch struct {
	Fields []string
	Args   []interface{}
	Index  int
}

// SQLPatches parse the resource to be patched in the db (only update)
// source : https://play.golang.org/p/TdwAhb7pjT
func SQLPatches(resource interface{}) SQLPatch {
	var sqlPatch SQLPatch
	rType := reflect.TypeOf(resource)
	rVal := reflect.ValueOf(resource)
	n := rType.NumField()

	sqlPatch.Fields = make([]string, 0, n)
	sqlPatch.Args = make([]interface{}, 0, n)
	sqlPatch.Index = 0

	for i := 0; i < n; i++ {
		fType := rType.Field(i)
		fVal := rVal.Field(i)
		tag := fType.Tag.Get("db")

		// skip nil properties (not going to be patched), skip unexported fields, skip fields to be skipped for SQL
		if fVal.IsNil() || fType.PkgPath != "" || tag == "-" || tag == "" {
			continue
		}

		// if no tag is set, use the field name
		if tag == "" {
			tag = fType.Name
		}
		// and make the tag lowercase in the end
		tag = strings.ToLower(tag)

		sqlPatch.Index++
		sqlPatch.Fields = append(sqlPatch.Fields, tag+" = $"+fmt.Sprintf("%d", sqlPatch.Index))

		var val reflect.Value
		if fVal.Kind() == reflect.Ptr {
			val = fVal.Elem()
		} else {
			val = fVal
		}

		switch val.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			sqlPatch.Args = append(sqlPatch.Args, val.Int())
		case reflect.String:
			sqlPatch.Args = append(sqlPatch.Args, val.String())
		case reflect.Bool:
			if val.Bool() {
				sqlPatch.Args = append(sqlPatch.Args, 1)
			} else {
				sqlPatch.Args = append(sqlPatch.Args, 0)
			}
		}
	}

	return sqlPatch
}

// GetUpdateQuery build and return a db query string
func (sqlPatch *SQLPatch) GetUpdateQuery(table string) string {
	query := "UPDATE " + table + " SET "
	for _, field := range sqlPatch.Fields {
		query = query + field + ","
	}
	return strings.TrimRight(query, ",")
}
