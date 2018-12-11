package errcode_def

import (
	"bufio"
	"strings"
)

type lineType int

const (
	moduleDefine lineType = iota
	clientVariableDefine
	serverVariableDefine
	otherDefine
)

type lineIterator interface {
	Next() bool
	Err() error
	Line() string
}

type lineContext interface {
	LineType() lineType
	Line() string
	Iterator() lineIterator
}

type scannerIterator struct {
	scanner *bufio.Scanner
	line    string
}

func newScannerIterator(scanner *bufio.Scanner) lineIterator {
	return &scannerIterator{scanner: scanner}
}

func (iter *scannerIterator) Next() bool {
	ok := iter.scanner.Scan()
	if ok {
		iter.line = trimRawLine(iter.scanner.Text())
	} else {
		iter.line = ""
	}
	return ok
}

func (iter *scannerIterator) Err() error {
	return iter.scanner.Err()
}

func (iter *scannerIterator) Line() string {
	return iter.line
}

type lineScannerContext struct {
	iter lineIterator
}

func newLineScannerContext(scanner *bufio.Scanner) lineContext {
	return &lineScannerContext{iter: newScannerIterator(scanner)}
}

func (c *lineScannerContext) Iterator() lineIterator {
	return c.iter
}

func (c *lineScannerContext) Line() string {
	return c.Line()
}

func (c *lineScannerContext) LineType() lineType {
	return parseLineType(c.Line())
}

func trimRawLine(line string) string {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
		return ""
	}
	idx := strings.IndexByte(line, '#')
	if idx != -1 {
		line = strings.TrimSpace(line[:idx])
	}
	return line
}

func parseLineType(line string) lineType {
	sz := len(line)
	if sz < 3 {
		// min lenght len([A])
		return otherDefine
	}
	if sz > 4 && line[:2] == "[[" && line[len(line)-2:] == "]]" {
		switch strings.TrimSpace(line[2 : len(line)-2]) {
		case "client_error":
			return clientVariableDefine
		case "server_error":
			return serverVariableDefine
		default:
			return otherDefine
		}
	}
	if line[0] == '[' && line[sz-1] == ']' && line[1] != '[' {
		return moduleDefine
	}
	return otherDefine
}
