//go:build windows
// +build windows

package main

import (
	"fmt"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"os"
	"path/filepath"
)

func exePath() (string, error) {
	//Getting the absolute path to the exe
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}

	//Checking whether the path is a directory
	//Returns an error if the path doesn't exist,
	//usually because the path lacks an .exe extension
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is a directory", p)
	}

	//Checking if the path lacks an .exe extension
	//If so, add an .exe extension to the path
	if filepath.Ext(p) == "" {
		p += ".exe"
		//Here we need to check whether it is a directory again
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is a directory", p)
		}
	}

	return "", err
}

func installService(name, desc string) error {
	//Getting the absolute path to the executable
	exepath, err := exePath()
	if err != nil {
		return err
	}

	//Connect to the Windows service manager
	mngr, err := mgr.Connect()
	if err != nil {
		return err
	}
	//Defer disconnect from the Windows service manager
	defer func(mngr *mgr.Mgr) {
		err := mngr.Disconnect()
		if err != nil {
			log.Println("Error disconnecting from Windows service manager: ", err)
		}
	}(mngr)

	//Checking whether the service is already installed
	s, err := mngr.OpenService(name)
	if err == nil {
		err = s.Close()
		if err != nil {
			log.Println("Failed to close connection to the service: ", err)
		}
		return fmt.Errorf("service %s already exists", name)
	}

	//Create new Windows service
	s, err = mngr.CreateService(name, exepath, mgr.Config{DisplayName: desc}, "is", "auto-started")
	if err != nil {
		return err
	}

	//Defer disconnect from the service
	defer func(s *mgr.Service) {
		err := s.Close()
		if err != nil {
			log.Println("Failed to close connection to the service: ", err)
		}
	}(s)

	return nil
}

func removeService(name string) error {
	//Connect to the Windows service manager
	mngr, err := mgr.Connect()
	if err != nil {
		return err
	}
	//Defer disconnect from the Windows service manager
	defer func(mngr *mgr.Mgr) {
		err := mngr.Disconnect()
		if err != nil {
			log.Println("Error disconnecting from Windows service manager: ", err)
		}
	}(mngr)

	// Getting the service to remove
	// Error signals that the service not installed
	s, err := mngr.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}

	//Defer disconnect from the service
	defer func(s *mgr.Service) {
		err := s.Close()
		if err != nil {
			log.Println("Failed to close connection to the service: ", err)
		}
	}(s)

	err = s.Delete()
	if err != nil {
		return err
	}
	return nil
}
