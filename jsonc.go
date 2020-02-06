package jsonc

import (
	"errors"
	"io"
)

type (
	state int

	comment struct {
		state   state
		multiLn bool
		isJSON  bool
	}

	// Decoder implements io.Reader. It wraps provided source.
	Decoder struct {
		r io.Reader
		c comment
	}
)

const (
	// byte representations of string literals
	tab      = 9   // (	)
	newLine  = 10  // (\n)
	space    = 32  // ( )
	quote    = 34  // ("")
	star     = 42  // (*)
	fwdSlash = 47  // (/)
	bkdSlash = 92  // (\)
	charN    = 110 // (n)
)

const (
	stopped  state = iota // 0, default state
	canStart              // 1
	started               // 2
	canStop               // 3
)

var (
	ErrUnexpectedEndOfJSON = errors.New("unexpected end of json")
)

// New a new io.Reader wrapping the provided one.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		c: comment{},
		r: r,
	}
}

// Read reads from underlying writer and processes the stream to omit comments.
// A single read doesn't guaranttee a valid JSON. Depends on length of passed slice.
//
// Produces ErrUnexpectedEndOfJSON for incomplete comments
func (d *Decoder) Read(p []byte) (int, error) {

	n, err := d.r.Read(p)
	if err != nil {
		return n, err
	}

	shortRead := n <= len(p)
	n = d.decode(p[:n])

	if shortRead && d.c.state != stopped {
		return 0, ErrUnexpectedEndOfJSON
	}

	return n, nil
}

func (d *Decoder) decode(p []byte) int {
	i := 0
	for _, s := range p {
		if d.c.handle(s) {
			p[i] = s
			i++
		}
	}

	return i
}

// handle the current byte, if returned true, add the byte to result.
func (c *comment) handle(s byte) bool {
	switch c.state {

	case stopped:
		if s == quote { // all characters between "" are valid, can be added to result
			c.isJSON = !c.isJSON
		}

		if c.isJSON {
			return true
		}

		if s == space || s == tab || s == newLine {
			return false
		}

		if s == fwdSlash {
			c.state = canStart
			return false
		}

		return true

	case canStart:

		if s == star || s == fwdSlash {
			c.state = started
		}

		c.multiLn = (s == star)

	case started:

		if s == star || s == bkdSlash {
			c.state = canStop
		}

		if s == newLine && !c.multiLn {
			c.state = stopped
		}

	case canStop:

		if s == fwdSlash || s == charN {
			c.state = stopped
			c.multiLn = false
		}

	}

	return false
}
