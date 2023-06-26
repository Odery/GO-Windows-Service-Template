//go:build windows

package main

import (
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"time"
)

func startService(name string) error {
	// Connect to the service manager
	mngr, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer func(mngr *mgr.Mgr) {
		err := mngr.Disconnect()
		if err != nil {
			log.Println("error disconnecting from the service manager")
		}
	}(mngr)

	// Open the service
	s, err := mngr.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access the service: %v", err)
	}
	defer func(s *mgr.Service) {
		err := s.Close()
		if err != nil {
			log.Println("error disconnecting from the service")
		}
	}(s)

	// Start the service
	err = s.Start("is", "manual-started")
	if err != nil {
		return fmt.Errorf("could not start the service: %v", err)
	}
	return nil
}

func controlService(name string, command svc.Cmd, to svc.State) error {
	// Connect to the service manager
	mngr, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer func(mngr *mgr.Mgr) {
		err := mngr.Disconnect()
		if err != nil {
			log.Println("error disconnecting from the service manager")
		}
	}(mngr)

	// Open the service
	s, err := mngr.OpenService(name)
	if err != nil {
		return fmt.Errorf("could not access the service: %v", err)
	}
	defer func(s *mgr.Service) {
		err := s.Close()
		if err != nil {
			log.Println("error disconnecting from the service")
		}
	}(s)

	// Control the service
	status, err := s.Control(command)
	if err != nil {
		return fmt.Errorf("could not send control command = '%d', error: %v", command, err)
	}

	// Polling service for changes
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout for the service to change state on command: %d", command)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retireve service status after command, error: %v", err)
		}
	}
	return nil
}
