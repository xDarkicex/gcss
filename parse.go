package gcss

import "strings"

// Special characters
const (
	cr   = "\r"
	lf   = "\n"
	crlf = "\r\n"
)

// parse parses the string, generates the elements
// and returns the two channels: the first one returns
// the generated elements and the last one returns
// an error when it occurs.
func parse(s string) (<-chan []element, <-chan error) {
	elemsc := make(chan []element)
	errc := make(chan error)

	go func() {
		var elems []element

		lines := strings.Split(formatLF(s), lf)

		i := 0
		l := len(lines)

		for i < l {
			// Fetch a line.
			ln := newLine(i+1, lines[i])
			i++

			// Ignore the empty line.
			if ln.isEmpty() {
				continue
			}

			if ln.isTopIndent() {
				e := newElement(ln, nil)

				if err := appendChildren(e, lines, &i, l); err != nil {
					errc <- err
					return
				}

				elems = append(elems, e)
			}
		}

		elemsc <- elems
	}()

	return elemsc, errc
}

// appendChildren parses the lines and appends the child elements
// to the parent element.
func appendChildren(parent element, lines []string, i *int, l int) error {
	for *i < l {
		// Fetch a line.
		ln := newLine(*i+1, lines[*i])

		// Ignore the empty line.
		if ln.isEmpty() {
			*i++
			return nil
		}

		ok, err := ln.childOf(parent)

		if err != nil {
			return err
		}

		if !ok {
			return nil
		}

		child := newElement(ln, parent)

		parent.AppendChild(child)

		*i++

		if err := appendChildren(child, lines, i, l); err != nil {
			return err
		}
	}

	return nil
}

// formatLF replaces the line feed codes with LF and
// returns the result string.
func formatLF(s string) string {
	return strings.Replace(strings.Replace(s, crlf, lf, -1), cr, lf, -1)
}
