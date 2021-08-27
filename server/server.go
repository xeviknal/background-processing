package server

import (
	"math/rand"
	"time"

	"github.com/xeviknal/background-commons/database"
)

type Server struct {
	StartedAt time.Time
}

func NewServer() *Server {
	return &Server{
		StartedAt: time.Now(),
	}
}

func (s *Server) Start() {
	// Starting a seed for randoms
	rand.Seed(time.Now().UnixNano())

	// Setting appropriate db connection
	database.SetConnectionConfig("jobs", "jobs", "jobs")

	go func() {
		for {
			time.Sleep(5 * time.Second)
		}
	}()
}

func (s *Server) Stop() {
	database.Close()
}
