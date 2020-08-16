package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "filename required\n")
		os.Exit(1)
	}

	config, err := ReadConfig(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	o, err := json.Marshal(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON error: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Println(string(o))
}
