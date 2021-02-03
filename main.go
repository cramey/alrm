package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"git.binarythought.com/cdramey/alrm/config"
	"os"
	"strings"
)

func main() {
	cfgPath := flag.String("c", "", "path to configuration file")
	debuglvl := flag.Int("d", 0, "debug level")

	flag.Usage = printUsage
	flag.Parse()

	if *cfgPath == "" {
		searchpaths := []string{"/etc/alrmrc", "./alrmrc"}
		for _, sp := range searchpaths {
			if _, err := os.Stat(sp); err == nil {
				*cfgPath = sp
				break
			}
		}
		if *cfgPath == "" {
			fmt.Fprintf(os.Stderr, "cannot find configuration\n")
			os.Exit(1)
		}
	}

	command := strings.ToLower(flag.Arg(0))
	switch command {
	case "config":
		if *debuglvl > 0 {
			fmt.Printf("checking config %s .. \n", *cfgPath)
		}

		cfg, err := config.ReadConfig(*cfgPath, *debuglvl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		if *debuglvl > 0 {
			o, err := json.MarshalIndent(cfg, "", "  ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "JSON error: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("%s\n", string(o))
		}

		fmt.Printf("config is OK\n")

	case "alarm":
		an := flag.Arg(1)
		if an == "" {
			fmt.Fprintf(os.Stderr, "alarm name required\n")
			os.Exit(1)
		}

		cfg, err := config.ReadConfig(*cfgPath, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		al, exists := cfg.Alarms[an]
		if !exists {
			fmt.Fprintf(os.Stderr, "group or host is not defined\n")
			os.Exit(1)
		}

		err = al.Alarm(
			"test group", "test host", "test check",
			fmt.Errorf("test alarm message"),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "alarm failed: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("alarm sounded successfully\n")

	case "check":
		cn := flag.Arg(1)
		if cn == "" {
			fmt.Fprintf(os.Stderr, "check host or group name required\n")
			os.Exit(1)
		}

		cfg, err := config.ReadConfig(*cfgPath, 0)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		gr, exists := cfg.Groups[cn]
		if !exists {
			fmt.Fprintf(os.Stderr, "group or host is not defined\n")
			os.Exit(1)
		}

		err = gr.Check(*debuglvl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "check failed: %s\n", err.Error())
			os.Exit(1)
		}
		fmt.Printf("check successful\n")

	case "server":
		cfg, err := config.ReadConfig(*cfgPath, *debuglvl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

		err = startServer(cfg, *debuglvl)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			os.Exit(1)
		}

	case "help", "":
		printUsage()

	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", command)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Printf("Usage: %s [args] <action> ...\n", os.Args[0])
	fmt.Printf("Arguments:\n")
	fmt.Printf("  -c <path>  : path to configuration file\n")
	fmt.Printf("  -d <level> : debug level (0-9, higher for more debugging)\n")
	fmt.Printf("Actions:\n")
	fmt.Printf("  verify configuration:     %s [args] config\n", os.Args[0])
	fmt.Printf("  run a check manually:     %s [args] check <host/group>\n", os.Args[0])
	fmt.Printf("  test an alarm:            %s [args] alarm <name>\n", os.Args[0])
	fmt.Printf("  start server (forground): %s [args] server\n", os.Args[0])
}
