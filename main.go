//go:build windows

package main

import (
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"log"
	"os"
	"strings"
)

// Default name value for the service
var (
	svcName = "TemplateService"
	svcDsc  = "TemplateService"
)

func main() {
	// Parsing command line arguments
	flag.StringVar(&svcName, "name", "TemplateService", "name of the service")
	flag.StringVar(&svcDsc, "description", "", "description of the service")
	flag.Parse()

	//Checking whether the app is running as a service
	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service mode: %v", err)
	}

	// Run the service if so
	if inService {
		runService(svcName)
		return
	}

	//Checking if arguments were specified
	if len(flag.Args()) < 1 {
		usage("no command specified")
	}

	//Checking whether the service name was specified
	if svcName == "TemplateService" {
		usage("service name was not specified")
	}

	// Parsing command
	cmd := strings.ToLower(flag.Args()[0])

	// Determining what command was specified
	switch cmd {
	case "install":
		err = installService(svcName, svcDsc)
	case "remove":
		err = removeService(svcName)
	case "start":
		err = startService(svcName)
	case "stop":
		err = controlService(svcName, svc.Stop, svc.Stopped)
	case "pause":
		err = controlService(svcName, svc.Pause, svc.Paused)
	case "continue":
		err = controlService(svcName, svc.Continue, svc.Running)
	default:
		usage(fmt.Sprintf("invalid command: %v", cmd))
	}

	if err != nil {
		log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
}

func usage(errmsg string) {
	_, err := fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	if err != nil {
		log.Println("Error printing usage message to the 'Stderr': ", err)
	}

	os.Exit(2)
}
