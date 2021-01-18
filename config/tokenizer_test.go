package config

import (
	"encoding/json"
	"testing"
)

func TestSimpleSpaces(t *testing.T) {
	runTest(t, "simple-spaces",
		`[["one","two","three","four","five","six"]]`,
	)
}

func TestSimpleMultiline(t *testing.T) {
	runTest(t, "simple-multiline",
		`[["one","two","three"],["four","five"],[],[],["six"]]`,
	)
}

func TestQuotes(t *testing.T) {
	runTest(t, "quotes",
		`[["one","two","three"],[],["four","five","six"]]`,
	)
}

func TestQuotesMultiline(t *testing.T) {
	runTest(t, "quotes-multiline",
		`[["one\ntwo"],["three\nfour"],[],[],["five\n  six"]]`,
	)
}

func TestQuotesEmpty(t *testing.T) {
	runTest(t, "quotes-empty",
		`[["one","","three"],["","five",""],["seven"]]`,
	)
}

func TestComments(t *testing.T) {
	runTest(t, "comments",
		`[[],["one"],[],["two"],[],["three"]]`,
	)
}

func TestCommentsInline(t *testing.T) {
	runTest(t, "comments-inline",
		`[["one"],["two#three"],[],["four"]]`,
	)
}

func runTest(t *testing.T, bn string, exp string) {
	t.Logf("Running testdata/%s.tok.. ", bn)
	tok, err := NewTokenizer("testdata/" + bn + ".tok")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer tok.Close()

	tokens := [][]string{}
	for tok.Scan() {
		ln := tok.Line()
		tl := len(tokens)
		if tl < ln {
			for i := tl; i < ln; i++ {
				tokens = append(tokens, []string{})
			}
		}
		tokens[ln-1] = append(tokens[ln-1], tok.Text())
	}
	if tok.Err() != nil {
		t.Fatalf("%s", tok.Err())
	}

	out, err := json.Marshal(tokens)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if exp != string(out) {
		t.Logf("Expected: %s", exp)
		t.Logf("Got:      %s", out)
		t.FailNow()
	}
}
