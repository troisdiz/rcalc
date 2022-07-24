package rcalc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func debugLex(name string, input string) *Lexer {
	lexer := Lex(name, input)
	lexer.debugMode = true
	return lexer
}

func lexForTest(txt string) []LexItem {
	lex := debugLex("Test", txt)

	var items []LexItem

	for it := lex.NextItem(); it.typ != lexItemEOF; it = lex.NextItem() {
		items = append(items, it)
	}
	return items
}

func TestLexerInteger(t *testing.T) {
	txt := " 123 7.9 456e29 -2 +3.5"
	items := lexForTest(txt)
	for _, it := range items {
		fmt.Printf("%v\n", it)
	}
	assert.Len(t, items, 5, "Wrong number of lexItems have been returned")
	assert.Equal(t, items[0].typ, lexItemNumber)
	assert.Equal(t, items[0].val, "123")
	assert.Equal(t, items[1].typ, lexItemNumber)
	assert.Equal(t, items[1].val, "7.9")
	assert.Equal(t, items[2].typ, lexItemNumber)
	assert.Equal(t, items[2].val, "456e29")
	assert.Equal(t, items[3].typ, lexItemNumber)
	assert.Equal(t, items[3].val, "-2")
	assert.Equal(t, items[4].typ, lexItemNumber)
	assert.Equal(t, items[4].val, "+3.5")
}

func TestLexerIdentifier(t *testing.T) {
	txt := " 'ab' "
	items := lexForTest(txt)
	for _, it := range items {
		fmt.Printf("%v\n", it)
	}
	if assert.Len(t, items, 1, "1 lexItem should be returned") {
		assert.Equal(t, items[0].typ, lexItemIdentifier)
		assert.Equal(t, items[0].val, "'ab'")
	}
}

func TestLexerActionKeyword(t *testing.T) {
	txt := " ab then "
	items := lexForTest(txt)
	for _, it := range items {
		fmt.Printf("%v\n", it)
	}
	assert.Len(t, items, 2, "1 lexItem should be returned")
}

func TestLexerList(t *testing.T) {
	txt := " {'ab'} "
	items := lexForTest(txt)
	fmt.Println("Lex Items: start")
	for _, it := range items {
		fmt.Printf("%v\n", it)
	}
	fmt.Println("Lex Items: end")

	if assert.Len(t, items, 3, "1 lexItem should be returned") {
		assert.Equal(t, items[0].typ, lexItemOpenCurlyBrace)
		assert.Equal(t, items[0].val, "{")
		assert.Equal(t, items[1].typ, lexItemIdentifier)
		assert.Equal(t, items[1].val, "'ab'")
		assert.Equal(t, items[2].typ, lexItemCloseCurlyBrace)
		assert.Equal(t, items[2].val, "}")
	}
}
