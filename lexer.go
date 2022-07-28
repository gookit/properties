package properties

import (
	"fmt"
	"strings"
	"text/scanner"

	"github.com/gookit/goutil/errorx"
)

type lexer struct {
	text string
}

func lex(text string) *lexer {
	return &lexer{
		text: text,
	}
}

func (l *lexer) parse() error {
	r := strings.NewReader(l.text)

	var s scanner.Scanner
	s.Init(r)
	// s.Mode = scanner.ScanIdents | scanner.ScanComments | scanner.ScanStrings | scanner.ScanRawStrings
	// s.Mode ^= scanner.SkipComments // don't skip comments. starts with: "/*" "//"
	s.Filename = "default"

	var line int
	var key, val, comments string
	var lineTok rune

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		// start new line
		isStart := line != s.Line
		if isStart {
			lineTok = tok
			fmt.Println("- newline")
		}
		line = s.Line

		switch lineTok {
		case '#', scanner.Comment: // comments line
			comments += s.TokenText()
		case scanner.Ident: // value line
			switch tok {
			case '.': // sep char
				if isStart {
					return errorx.Rawf("sep char cannot at line start, pos: %s", s.Position)
				}
			case '=':
				if key == "" {
					return errorx.Rawf("char '=' cannot at line start, pos: %s", s.Position)
				}
				val += s.TokenText()
			case scanner.Ident:
				if isStart {
					key += s.TokenText()
				}
			}
		}

		fmt.Printf(
			"%s: tok=%s txt=%s\n",
			s.Position, scanner.TokenString(tok), s.TokenText(),
		)
	}

	return nil
}

func (l *lexer) next() tokenItem {
	return tokenItem{}
}
