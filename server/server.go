package server

import (
	"math/rand"
	"time"

	"github.com/xeviknal/background-commons/database"

	"github.com/xeviknal/background-processing/consumers"
	"github.com/xeviknal/background-processing/publishers"
)

const defaultPublishCadence = 3 * time.Second

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
	database.SetConnectionConfig("jobs", "jobs", "jobs", "mariadb.mariadb.svc.cluster.local")

	go func() {
		go consumers.ConsumeTasks()
		for {
			go publishers.PublishTasks()
			time.Sleep(defaultPublishCadence)
		}
	}()
}

func (s *Server) Stop() {
	database.Close()
}
