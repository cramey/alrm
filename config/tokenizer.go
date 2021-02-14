package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const (
	TK_NONE = iota
	TK_VAL
	TK_QUOTE
	TK_COMMENT
)

type Tokenizer struct {
	curline int
	repline int
	file    *os.File
	reader  *bufio.Reader
	text    string
	err     error
}

func NewTokenizer(fn string) (*Tokenizer, error) {
	var err error
	tk := &Tokenizer{curline: 1}
	tk.file, err = os.Open(fn)
	if err != nil {
		return nil, err
	}

	tk.reader = bufio.NewReader(tk.file)
	return tk, nil
}

func (t *Tokenizer) Close() error {
	return t.file.Close()
}

func (t *Tokenizer) Scan() bool {
	t.repline = t.curline
	state := TK_NONE
	t.text = ""

	var b strings.Builder
	var quo rune
	for {
		var r rune
		r, _, t.err = t.reader.ReadRune()
		if t.err != nil {
			break
		}
		if r == unicode.ReplacementChar {
			t.err = fmt.Errorf("invalid utf-8 encoding on line %s", t.repline)
			break
		}

		switch state {
		case TK_NONE:
			// When between values, increment both the reported line
			// and the current line, since there's not yet anything
			// to report
			if r == '\n' {
				t.repline++
				t.curline++
			}

			// If we're between values and we encounter a space
			// or a control character, ignore it
			if unicode.IsSpace(r) || unicode.IsControl(r) {
				continue
			}

			// If we're between values and we encounter a #, it's
			// the beginning of a comment
			if r == '#' {
				state = TK_COMMENT
				continue
			}

			// If we're between values and we get a quote character
			// treat it as the beginning of a string literal
			if r == '"' || r == '\'' || r == '`' {
				state = TK_QUOTE
				quo = r
				continue
			}

			b.WriteRune(r)
			state = TK_VAL

		case TK_VAL:
			// In values, only increment the current line, so
			// if an error is reported, it reports the line
			// the value starts on
			if r == '\n' {
				t.curline++
			}

			// If we're in a normal value and we encounter a space
			// or a control character, end value
			if unicode.IsSpace(r) || unicode.IsControl(r) {
				goto end
			}
			b.WriteRune(r)

		case TK_QUOTE:
			// In quotes, only increment the current line, so
			// if an error is reported, it reports the line
			// the quoted value starts on
			if r == '\n' {
				t.curline++
			}

			// End this quote if it's another quote of the same rune
			if r == quo {
				goto end
			}
			b.WriteRune(r)

		case TK_COMMENT:
			// Comments are ignored, until a new line is encounter
			// at which point, increment the current and reported line
			if r == '\n' {
				t.curline++
				t.repline++
				state = TK_NONE
			}
			continue
		}
	}

end:
	if t.err == nil || t.err == io.EOF {
		t.text = b.String()
	}
	return t.err == nil
}

func (t *Tokenizer) Text() string {
	return t.text
}

func (t *Tokenizer) Line() int {
	return t.repline
}

func (t *Tokenizer) Err() error {
	if t.err == io.EOF {
		return nil
	}
	return t.err
}
