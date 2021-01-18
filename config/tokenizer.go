package config

import (
	"bufio"
	"fmt"
	"os"
)

type Tokenizer struct {
	line    int
	file    *os.File
	scanner *bufio.Scanner
}

func NewTokenizer(fn string) (*Tokenizer, error) {
	var err error
	tk := &Tokenizer{line: 1}
	tk.file, err = os.Open(fn)
	if err != nil {
		return nil, err
	}

	tk.scanner = bufio.NewScanner(tk.file)
	tk.scanner.Split(tk.Split)
	return tk, nil
}

func (t *Tokenizer) Close() error {
	return t.file.Close()
}

func (t *Tokenizer) Scan() bool {
	return t.scanner.Scan()
}

func (t *Tokenizer) Text() string {
	return t.scanner.Text()
}

func (t *Tokenizer) Line() int {
	return t.line
}

func (t *Tokenizer) Err() error {
	return t.scanner.Err()
}

func (t *Tokenizer) Split(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var ignoreline bool
	var started bool
	var startidx int
	var quote byte

	for i := 0; i < len(data); i++ {
		c := data[i]
		//fmt.Printf("%c (%t) (%t)\n", c, started, ignoreline)
		switch c {
		case '\f', '\n', '\r':
			if started {
				return i, data[startidx:i], nil
			}

			t.line++
			if ignoreline {
				ignoreline = false
				continue
			}
			fallthrough

		case ' ', '\t', '\v':
			if started && quote == 0 {
				return i + 1, data[startidx:i], nil
			}

		case '\'', '"', '`':
			// When the quote ends
			if quote == c {
				// if we've gotten data, return it
				if started {
					return i + 1, data[startidx:i], nil
				}
				// if we haven't return nothing
				return i + 1, []byte{}, nil
			}

			// start a quoted string
			if !ignoreline && quote == 0 {
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
		if quote != 0 {
			return 0, nil, fmt.Errorf("unterminated quote")
		}

		if ignoreline {
			return len(data), nil, nil
		}
		if started {
			return len(data), data[startidx:], nil
		}
	}

	return 0, nil, nil
}
