package agents

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	StatusZero = iota
	StatusInDelimiter

	LineTypeLine        string = "line"
	LineTypeBODelimiter        = "bo-delimiter"
	LineTypeEODelimiter        = "eo-delimiter"
)

func ExtractTextForParts(text string, delimiters []string) (map[string]string, error) {
	const semLogContext = semLogPackageContext + "extract-xml-sections-from-text"

	var beginDelimiters []string
	var endDelimiters []string
	for _, d := range delimiters {
		beginDelimiters = append(beginDelimiters, fmt.Sprintf("<%s>", d))
		endDelimiters = append(endDelimiters, fmt.Sprintf("</%s>", d))
	}

	bytesReader := strings.NewReader(text)
	bufReader := bufio.NewReader(bytesReader)

	var m map[string]string
	var sb strings.Builder
	var currentDelimiter, delim string
	var lineType string
	status := StatusZero
	line, isPrefix, err := bufReader.ReadLine()
	for err == nil && !isPrefix {

		sline := string(line)
		sline, lineType, delim = typeOfLine(sline, delimiters)

		switch lineType {
		case LineTypeBODelimiter:
			if status != StatusZero {
				return nil, fmt.Errorf("found %s when status in %d (%s)", delim, status, currentDelimiter)
			}
			currentDelimiter = delim
			if sline != "" {
				sb.WriteString(sline)
				sb.WriteString("\n")
			}
			status = StatusInDelimiter
		case LineTypeEODelimiter:
			if status != StatusInDelimiter || (status == StatusInDelimiter && currentDelimiter != delim) {
				return nil, fmt.Errorf("found end of delimiter %s when status in %d and current delimiter is %s", delim, status, currentDelimiter)
			}
			if sline != "" {
				sb.WriteString(sline)
				sb.WriteString("\n")
			}
			status = StatusZero
			if m == nil {
				m = make(map[string]string)
			}
			m[currentDelimiter] = sb.String()
			sb = strings.Builder{}
		default:
			if status != StatusZero {
				sb.WriteString(sline)
				sb.WriteString("\n")
			}
		}

		line, isPrefix, err = bufReader.ReadLine()
	}

	if err != nil && err != io.EOF {
		return nil, err
	}

	if isPrefix {
		return nil, errors.New("buffer too small")
	}

	return m, nil
}

func typeOfLine(line string, delimiters []string) (string, string, string) {
	const semLogContext = semLogPackageContext + "type-of-line"

	// NOTE: delimiter detection is done on a trimmed copy so that indented
	// opening/closing tags are still recognized, but plain content lines are
	// returned verbatim (untrimmed) to preserve their indentation. Only the
	// inline remainder that shares a line with a tag is trimmed.
	trimmed := strings.TrimSpace(line)
	for _, d := range delimiters {
		start := fmt.Sprintf("<%s>", d)
		if ndx := strings.Index(trimmed, start); ndx >= 0 {
			return strings.TrimSpace(trimmed[ndx+len(start):]), LineTypeBODelimiter, d
		}

		end := fmt.Sprintf("</%s>", d)
		if ndx := strings.Index(trimmed, end); ndx >= 0 {
			return strings.TrimSpace(trimmed[:ndx]), LineTypeEODelimiter, d
		}
	}

	return line, LineTypeLine, ""
}
