// package jsonc provides APIs to read commented (.jsonc) files.
// It implements io.Reader which makes it compatible with standard library.

// EX.
// 		...
// 		ba, _ := ioutil.Readall(jsonc.NewDecoder(r)) // r is io.Reader
// 		...

package jsonc

import (
	"errors"
	"io"
	"reflect"
	"unsafe"
)

type (
	// Decoder represents a .jsonc decoder. It implements io.Reader.
	Decoder struct {
		r io.Reader
		c comment
	}

	// comment represents status of Decoder while reading json document.
	comment struct {
		state   state
		multiLn bool // if true, a multiline comment
		isJSON  bool // if true, decoder is processing either a "key" or "value"
	}

	state int
)

const (
	stopped  state = iota // 0, default state, started with 0 to make `zero value` of `comment` useful.
	canStart              // 1
	started               // 2
	canStop               // 3
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

var (
	ErrUnexpectedEndOfComment = errors.New("unexpected end of comment")
)

func (c *comment) reset() {
	c.state = stopped
	c.multiLn = false
}

func (c *comment) complete() bool {
	return c.state == stopped
}

func decode(p []byte, c *comment) int {
	i := 0
	for _, s := range p {

		switch c.state {

		case stopped:
			if s == quote { // all characters between "" are valid, can be added to result
				c.isJSON = !c.isJSON
			}

			if c.isJSON {
				p[i] = s
				i++
				continue
			}

			// TODO: formatting can be preserved by allowing space and tab in JSON
			if s == space || s == tab || s == newLine {
				continue
			}

			if s == fwdSlash {
				c.state = canStart
				continue
			}

			p[i] = s
			i++

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
				c.reset()
			}

		case canStop:

			if s == fwdSlash || s == charN {
				c.reset()
			}

		}

	}

	return i
}

// NewDecoder returns a new Decoder wrapping the provided io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		c: comment{},
		r: r,
	}
}

// Read reads from underlying reader and processes the stream to omit comments.
// A single read doesn't guaranttee a valid JSON. Depends on length of passed slice.
//
// Produces ErrUnexpectedEndOfComment for incomplete comments else if not nil, then it is an error from underlying reader.
func (d *Decoder) Read(p []byte) (int, error) {

	n, err := d.r.Read(p)
	if err != nil {
		return 0, err
	}

	shortRead := n <= len(p) // data to be processed is less than length of `p`
	n = decode(p[:n], &d.c)

	if shortRead && !d.c.complete() {
		return 0, ErrUnexpectedEndOfComment
	}

	return n, nil
}

// DecodeBytes decodes passed commented json byte slice to normal json.
// It modifies the passed slice. The passed slice must be refferred till returned count, if there is no error.
//
// The error doesn't include errors related to invalid json. If not nil, it must be ErrUnexpectedEndOfComment.
//
// The returned json must be checked for validity.
func DecodeBytes(p []byte) (int, error) {
	c := &comment{}
	n := decode(p, c)

	if !c.complete() {
		return 0, ErrUnexpectedEndOfComment
	}

	return n, nil
}

// DecodeString decodes passed commented json to normal json.
// It uses "unsafe" way to convert a byte slice to result string. This saves allocations and improves performance is case of large json.
//
// The error doesn't include errors related to invalid json. If not nil, it must be ErrUnexpectedEndOfComment.
//
// The returned json must be checked for validity.
func DecodeString(s string) (string, error) {
	p := []byte(s)

	n, err := DecodeBytes(p)
	if err != nil {
		return "", err
	}

	p = p[:n]

	// following operation is safe to do till p is not being changed. This reduces allocations.
	sh := *(*reflect.SliceHeader)(unsafe.Pointer(&p))
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{
		Data: sh.Data,
		Len:  sh.Len,
	})), nil
}
