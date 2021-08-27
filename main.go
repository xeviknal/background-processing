package background_processing

import (
	"log"
	"os"
	"os/signal"

	"github.com/xeviknal/background-processing/server"
)

// Infinite loop taking background jobs
func main() {
	server := server.NewServer()
	server.Start()
	waitForTermination()
	server.Stop()
}

// Waiting until the process receive a Termination signal
func waitForTermination() {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	log.Println("Termination Signal Received")
}
