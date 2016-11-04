package lib

import (
  "encoding/json"
  "fmt"
  "strconv"
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
	var host, dbName, port, username, passwd =
		viper.GetString("db_host"),
		viper.GetString("db_name"),
		viper.GetString("db_port"),
		viper.GetString("db_username"),
		viper.GetString("db_password")

	var dbUrl = "postgres://"

	switch {
		case username != "":
			dbUrl += username
			dbUrl += ":" + passwd + "@"
			fallthrough
		case host != "":
			dbUrl += host
			fallthrough
		case port != "":
			dbUrl += ":" + port
			fallthrough
		case dbName != "":
			dbUrl += "/" + dbName
	}
	dbUrl += "?sslmode=disable"

	return dbUrl
}

type Query struct {
  Q    string
  Opt  *Option
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
