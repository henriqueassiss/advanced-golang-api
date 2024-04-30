package server

import (
	"log"

	"github.com/jwalton/gchalk"
)

func (s *Server) InitDomains() {
	log.Println(gchalk.Yellow("Domain: starting"))

	log.Println(gchalk.Blue("Domain: done"))

	if s.cfg.App.Environment == "development" {
		s.InitMocks()
	}
}

func (s *Server) InitMocks() {
	log.Println(gchalk.Yellow("Mock: starting"))

	log.Println(gchalk.Blue("Mock: done"))
}
