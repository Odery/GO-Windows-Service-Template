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
	flag.StringVar(&svcName, "name", svcName, "name of the service")
	flag.StringVar(&svcDsc, "description", svcDsc, "description of the service")
	flag.Parse()

	//Checking whether the app is running as a service
	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service mode: %v", err)
	}

	// Run the service if so
	if inService {
		//TODO runService(svcName, false)
		return
	}

	//Checking if arguments were specified
	if len(os.Args) < 2 {
		usage("no command specified")
	}

	//Checking whether the service name was specified
	if svcName == "TemplateService" {
		usage("service name was not specified")
	}

	// Parsing command
	cmd := strings.ToLower(os.Args[1])

	// Determining what command was specified
	switch cmd {
	case "debug":
		// TODO runService(svcName, true)
		return
	case "install":
		err = installService(svcName, svcDsc)
	case "remove":
		err = removeService(svcName)
	case "start":
		//TODO startService(svcName)
		return
	case "stop":
		//TODO controlService(svcName, svc.Stop, svc.Stopped)
		return
	case "pause":
		//TODO controlService(svcName, svc.Pause, svc.Paused)
		return
	case "continue":
		//TODO controlService(svcName, svc.Continue, svc.Running)
		return
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
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	if err != nil {
		log.Println("Error printing usage message to the 'Stderr': ", err)
	}

	os.Exit(2)
}
