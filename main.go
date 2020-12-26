package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"os"
	"strings"
)

func main() {
	configPath := flag.String("config", "", "path to configuration file")

	flag.Parse()

	if *configPath == "" {
		if _, err := os.Stat("./alrmrc"); err == nil {
			*configPath = "./alrmrc"
		}
		if _, err := os.Stat("/etc/alrmrc"); err == nil {
			*configPath = "/etc/alrmrc"
		}
		if *configPath == "" {
			fmt.Fprintf(os.Stderr, "Cannot find configuration\n")
			os.Exit(1)
		}
	}

	config, err := ReadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	command := strings.ToLower(flag.Arg(0))
	switch command {
		case "json":
			o, err := json.MarshalIndent(config, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "JSON error: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s", string(o))

		case "", "config":
			fmt.Fprintf(os.Stdout, "Config is OK.\n")
			os.Exit(0)

		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
			os.Exit(1)
	}
}
