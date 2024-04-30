package main

import (
	"github.com/henriqueassiss/advanced-golang-api/internal/server"
)

var Version = "v0.1.0"

func main() {
	s := server.New(server.WithVersion(Version))
	s.Init()
	s.Run()
}
