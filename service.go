//go:build windows

package main

import (
	"golang.org/x/sys/windows/svc"
	"log"
	"time"
)

type Service struct{}

func (s *Service) Execute(_ []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, erno uint32) {
	// Declare accepted commands
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue

	// Wait for the start
	changes <- svc.Status{State: svc.StartPending}

	// Declaring fastTick and slowTick
	fastTick := time.Tick(500 * time.Millisecond)
	slowTick := time.Tick(2 * time.Second)

	// Setting tick to the fast tick as default
	tick := fastTick

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case <-tick:
			// Do sth here
			log.Println("Simulating work...!")
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				// Sleep to avoid deadlock
				time.Sleep(100 * time.Millisecond)
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				// Exit on Stop ot Shutdown
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
				// Changing tick to a slow tick on pause
				tick = slowTick
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				// Changing tick back to the fast tick
				tick = fastTick
			default:
				log.Printf("Unexpected control request: %d\n", c)
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string) {
	log.Printf("Starting the service: %s\n", name)
	err := svc.Run(name, &Service{})

	if err != nil {
		log.Printf("Failed to start service: %s, error: %v\n", name, err)
		return
	}
	log.Printf("Service %s has stopped.", name)
}
