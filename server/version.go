package server

import (
  "fmt"
)

var (
  Version string
)

func PrintVersion() {
  fmt.Println(Version)
}
