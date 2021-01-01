package main

import (
	"encoding/json"
	"fmt"
	"flag"
	"os"
	"strings"
	"alrm/config"
)

func main() {
	cfgPath := flag.String("c", "", "path to configuration file")
	debuglvl := flag.Int("d", 0, "debug level")

	flag.Parse()

	if *cfgPath == "" {
		if _, err := os.Stat("./alrmrc"); err == nil {
			*cfgPath = "./alrmrc"
		}
		if _, err := os.Stat("/etc/alrmrc"); err == nil {
			*cfgPath = "/etc/alrmrc"
		}
		if *cfgPath == "" {
			fmt.Fprintf(os.Stderr, "Cannot find configuration\n")
			os.Exit(1)
		}
	}

	command := strings.ToLower(flag.Arg(0))
	switch command {
		case "json":
			cfg, err := config.ReadConfig(*cfgPath, *debuglvl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			o, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "JSON error: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Fprintf(os.Stdout, "%s\n", string(o))

		case "", "config":
			_, err := config.ReadConfig(*cfgPath, *debuglvl)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			fmt.Fprintf(os.Stdout, "Config is OK.\n")
			os.Exit(0)

		default:
			fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
			os.Exit(1)
	}
}
