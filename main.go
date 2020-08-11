package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"unicode/utf8"
)

func main(){
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

	lscan := bufio.NewScanner(file)
	lscan.Split(bufio.ScanLines)
	for lscan.Scan() {
		line := lscan.Text()
		// Ignore comments
		if len(line) < 1 || line[0] == '#' {
			continue
		}

		wscan := bufio.NewScanner(strings.NewReader(line))
		wscan.Split(ScanWords)
		for wscan.Scan() {
			word := wscan.Text()
			fmt.Printf("[%s] ", word)
		}
	}
	fmt.Printf("\n")
}

func ScanWords(data []byte, atEOF bool) (int, []byte, error) {
	start := 0
	quote := int32(0)
	for start < len(data) {
		r, w := utf8.DecodeRune(data[start:])
		if !isSpace(r) {
			if isQuote(r) {
				quote = r
				start += w
			}
			break
		}
		start += w
	}

	for i := start; i < len(data); {
		r, w := utf8.DecodeRune(data[i:])

		if (quote > 0 && quote == r) || (quote == 0 && isSpace(r)) {
			return i + w, data[start:i], nil
		}

		i += w
	}

	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	return start, nil, nil
}

func isQuote(r rune) bool {
	switch r {
		case '\u0022', '\u0027':
			return true
	}
	return false
}

func isSpace(r rune) bool {
	if r <= '\u00FF' {
		switch r {
			case ' ', '\t', '\n', '\v', '\f', '\r':
				return true
			case '\u0085', '\u00A0':
				return true
		}
		return false
	}

	if '\u2000' <= r && r <= '\u200a' {
		return true
	}

	switch r {
		case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}
