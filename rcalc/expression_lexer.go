package rcalc

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type LexItemType int

const (
	lexItemError LexItemType = iota // error occurred;
	// value is text of error
	lexItemEOF
	// itemBlank
	lexItemNumber
	lexItemIdentifier
	lexItemKeywword
	lexItemAction
)

const (
	quoteMeta rune = '\''
)

const eof = -1

var itemTypeByNames = []string{
	"lexItemError",
	"lexItemEOF",
	"itemBlank",
	"lexItemNumber",
	"lexItemIdentifier",
	"lexItemKeywword",
	"lexItemAction",
}

var keyWords = map[string]bool{
	"if":   true,
	"then": true,
	"else": true,
}

type LexItem struct {
	typ LexItemType // Type, such as lexItemNumber.
	val string      // Value, such as "23.2".
}

func (it LexItem) String() string {

	return fmt.Sprintf("[%s, val: %s]", itemTypeByNames[it.typ], it.val)
}

// Lexer holds the state of the scanner.
type Lexer struct {
	name      string       // used only for error reports.
	input     string       // the string being scanned.
	state     stateFn      // current state
	start     int          // start position of this LexItem.
	pos       int          // current position in the input.
	width     int          // width of last rune read from input.
	items     chan LexItem // channel of scanned items.
	debugMode bool         // debug mode print debug messages
}

func (l *Lexer) emit(t LexItemType) {
	l.items <- LexItem{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

// next returns the next rune in the input.
func (l *Lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	if l.debugMode {
		fmt.Printf("Next char is #%s#\n", string(r))
	}
	return r
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *Lexer) backup() {
	if l.debugMode {
		fmt.Printf("Backup\n")
	}
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) peekAndTestOrEOF(valid string) bool {
	r := l.peek()
	return (r == eof) || strings.ContainsRune(valid, r)
}

// accept consumes the next rune
// if it's from the valid set.
func (l *Lexer) accept(valid string) bool {
	return l.acceptFunc(func(r rune) bool {
		return strings.ContainsRune(valid, r)
	})
}

func (l *Lexer) acceptFunc(fn validFn) bool {
	if fn(l.next()) {
		return true
	} else {
		l.backup()
	}
	return false
}

/*
// acceptRun consumes a run of runes from the valid set.
func (l *Lexer) acceptRun(valid string) {
	for strings.ContainsRune(valid, l.next()) {
	}
	l.backup()
}

func (l *Lexer) acceptTestOrEnd(valid string) bool {
	if l.peek() == eof {
		return true
	} else {
		return l.accept(valid)
	}
}
*/

//validFn
type validFn func(r rune) bool

// acceptRunFunc consumes a run of runes while the provided validFn returns true
func (l *Lexer) acceptRunFunc(fn validFn) {
	for fn(l.next()) {
	}
	l.backup()
}

// stateFn represents the state of the scanner
// as a function that returns the next state.
type stateFn func(*Lexer) stateFn

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *Lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- LexItem{
		lexItemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

// NextItem returns the next LexItem from the input.
func (l *Lexer) NextItem() LexItem {

	for {
		select {
		case item := <-l.items:
			return item
		default:
			l.state = l.state(l)
		}
	}
	// panic("not reached")
}

// Lex creates a new scanner for the input string.
func Lex(name string, input string) *Lexer {
	l := &Lexer{
		name:  name,
		input: input,
		state: lexBlank,
		items: make(chan LexItem, 2), // Two items sufficient.
	}
	return l
}

func isNumericValidFn(r rune) bool {
	return unicode.IsDigit(r)
}

func isAlphaNumericValidFn(r rune) bool {
	return unicode.IsDigit(r) || unicode.IsLetter(r)
}

func isActionValidFn(r rune) bool {
	if unicode.IsDigit(r) || unicode.IsLetter(r) {
		return true
	} else {
		switch r {
		case '+' | '-' | '*' | '/':
			return true
		}
	}
	return false
}

func lexBlank(l *Lexer) stateFn {

	for {
		next := l.next()
		// fmt.Printf("Next is #%s#\n", string(next))
		switch {
		case next == ' ':
			l.ignore()
			continue
		case next == quoteMeta:
			return lexIdentifier
		// + and - => can be number or action
		case next == '+' || next == '-':
			return lexStartWithPlusMinus
		case unicode.IsDigit(next):
			return lexNumber
		case unicode.IsLetter(next):
			return lexAction
		}
		if next == eof {
			break
		}
	}
	l.emit(lexItemEOF) // Useful to make EOF a token.
	return nil         // Stop the run loop.
}

func lexStartWithPlusMinus(l *Lexer) stateFn {
	if unicode.IsDigit(l.peek()) {
		return lexNumber
	} else if l.peekAndTestOrEOF(" ") {
		l.emit(lexItemAction)
		return lexBlank
	} else {
		return l.errorf("+ or - without a number")
	}
}

func (l *Lexer) scanNumber() bool {
	// Optional leading sign.
	l.accept("+-")
	l.acceptRunFunc(isNumericValidFn)
	if l.accept(".") {
		l.acceptRunFunc(isNumericValidFn)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRunFunc(isNumericValidFn)
	}
	// Next thing mustn't be alphanumeric.
	if isAlphaNumericValidFn(l.peek()) {
		l.next()
		return false
	}
	return true
}

func lexNumber(l *Lexer) stateFn {
	if !l.scanNumber() {
		l.errorf("Bad number format %q", l.input[l.start:l.pos])
	}
	l.emit(lexItemNumber)
	fmt.Println("lexNumber -> lexBlank")
	return lexBlank
}

func lexIdentifier(l *Lexer) stateFn {

Loop:
	for {
		r := l.next()
		switch {
		case r == eof || r == '\n':
			return l.errorf("Unterminated identifier")

		case r == '\'':
			break Loop
		}
	}
	l.emit(lexItemIdentifier)
	if l.debugMode {
		fmt.Println("lexIdentifier -> lexBlank")
	}
	return lexBlank
}

func lexAction(l *Lexer) stateFn {
	l.acceptRunFunc(isActionValidFn)
	if l.accept(" ") || l.pos == len(l.input) {
		l.backup()
		itemStr := l.input[l.start:l.pos]
		if _, ok := keyWords[itemStr]; ok {
			l.emit(lexItemKeywword)
		} else {
			l.emit(lexItemAction)
		}
	}

	return lexBlank
}
