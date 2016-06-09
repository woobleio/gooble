package config

import (
  "fmt"
  "encoding/json"
  "os"
)

const (
  DBUser = ""
  DBPassword = ""
)

type Credentials struct {
  Username string
  Password string
}

func DBCredentials() Credentials {
  file, ferr := os.Open("/config/.config.json")
  creds := Credentials{}
  if ferr != nil {
    fmt.Print(ferr)
    return creds
  }
  decoder := json.NewDecoder(file)
  err := decoder.Decode(&creds)
  fmt.Print(err)
  return creds
}
