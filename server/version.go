package server

import (
  "fmt"
)

var (
  Version string
  Commit string
)

func PrintVersion() {
  fmt.Printf("%s-%s", Version, Commit)
}
