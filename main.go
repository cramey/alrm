package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "filename required\n")
		os.Exit(1)
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot open %s: %s\n",
			os.Args[1], err.Error())
		os.Exit(1)
	}
	defer file.Close()

	scan := bufio.NewScanner(file)
	scan.Split(Split)
	for scan.Scan() {
		word := scan.Text()
		fmt.Printf("[%s] ", word)
	}
	fmt.Printf("\n")
}

func Split(data []byte, atEOF bool) (int, []byte, error) {
	var ignoreline bool
	var started bool
	var startidx int
	var quote byte

	for i := 0; i < len(data); i++ {
		c := data[i]
		switch c {
		case '\f', '\n', '\r':
			if ignoreline {
				return i + 1, nil, nil
			}
			fallthrough

		case ' ', '\t', '\v':
			if started && quote == 0 {
				return i + 1, data[startidx:i], nil
			}

		case '\'', '"', '`':
			if started && quote == c {
				return i + 1, data[startidx:i], nil
			}

			if quote == 0 {
				quote = c
			}

		case '#':
			if !started {
				ignoreline = true
			}

		default:
			if !ignoreline && !started {
				started = true
				startidx = i
			}
		}
	}

	if atEOF {
		if ignoreline {
			return len(data), nil, nil
		}
		if started {
			return len(data), data[startidx:], nil
		}
	}

	return 0, nil, nil
}
