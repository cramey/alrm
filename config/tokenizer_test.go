package config

import (
	"testing"
	"encoding/json"
)

func TestTokenizer(t *testing.T) {
	runTest(t, "simple",
		`[["one","two","three","four","five","six"]]`,
	)
	runTest(t, "simple-broken",
		`[["one","two","three"],["four","five"],[],[],["six"]]`,
	)
	runTest(t, "comments",
	 `[[],["one","two","three"],[],["four","five","six"]]`,
	)
	runTest(t, "quotes",
	 `[["one","two three",[],["four five"],[],[" #six","","seven","ei","ght"],[],["multi\nline"]]`,
	)
}

func runTest(t *testing.T, bn string, exp string) {
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
		t.Logf("Got: %s", out)
		t.Fail()
	}
}
