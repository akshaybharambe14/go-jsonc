package jsonc

import (
	"bytes"
	"io"
	"reflect"
	"testing"
	"testing/iotest"
)

func ts(b []byte) *Decoder    { return &Decoder{r: bytes.NewBuffer(b)} }
func tsErr(b []byte) *Decoder { return &Decoder{r: iotest.DataErrReader(bytes.NewBuffer(b))} }

var (
	validSingle   = []byte(`{"foo": // this is a single line comment\n"bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com" }`)
	invalidSingle = []byte(`{"foo": // this is a single line comment "bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com" }`)

	validSingleESC   = []byte("{\"foo\": // this is a single line comment\n\"bar foo\", \"true\": false, \"number\": 42, \"object\": { \"test\": \"done\" }, \"array\" : [1, 2, 3], \"url\" : \"https://github.com\" }")
	invalidSingleESC = []byte("{\"foo\": // this is a single line comment\"bar foo\", \"true\": false, \"number\": 42, \"object\": { \"test\": \"done\" }, \"array\" : [1, 2, 3], \"url\" : \"https://github.com\" }")

	validBlock   = []byte(`{"foo": /* this is a block comment */ "bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com" }`)
	invalidBlock = []byte(`{"foo": /* this is a block comment "bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com" }`)
)

func Test_Decoder_Read(t *testing.T) {

	type args struct {
		p []byte
	}

	tests := []struct {
		name    string
		d       *Decoder
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Valid single line comment",
			d:       ts(validSingle),
			args:    args{p: make([]byte, len(validSingle))},
			want:    110, // (163(total) - 34(comments) - 19(spaces))
			wantErr: false,
		},
		{
			name:    "Invalid single line comment",
			d:       ts(invalidSingle),
			args:    args{p: make([]byte, len(invalidSingle))},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Valid single line comment (escaped json)",
			d:       ts(validSingleESC),
			args:    args{p: make([]byte, len(validSingleESC))},
			want:    110, // (163(total) - 34(comments) - 19(spaces))
			wantErr: false,
		},
		{
			name:    "Invalid single line comment (escaped json)",
			d:       ts(invalidSingleESC),
			args:    args{p: make([]byte, len(invalidSingleESC))},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Valid block comment",
			d:       ts(validBlock),
			args:    args{p: make([]byte, len(validBlock))},
			want:    110, // (159(total) - 29(comments) - 20(spaces))
			wantErr: false,
		},
		{
			name:    "Invalid block comment",
			d:       ts(invalidBlock),
			args:    args{p: make([]byte, len(invalidBlock))},
			want:    0,
			wantErr: true,
		},
		{
			name:    "Invalid Read",
			d:       tsErr(validBlock),
			args:    args{p: make([]byte, len(validBlock))},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Read(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decoder.Read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decoder.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewDecoder(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name string
		args args
		want *Decoder
	}{
		{
			name: "Valid Decoder",
			args: args{r: nil},
			want: &Decoder{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDecoder(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDecoder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeBytes(t *testing.T) {
	type args struct {
		p []byte
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "Valid input",
			args:    args{p: []byte(string(validBlock))},
			want:    110,
			wantErr: false,
		},
		{
			name:    "Invalid input",
			args:    args{p: []byte(string(invalidBlock))},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeBytes(tt.args.p)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Valid input",
			args:    args{s: string(validBlock)},
			want:    `{"foo":"bar foo","true":false,"number":42,"object":{"test":"done"},"array":[1,2,3],"url":"https://github.com"}`,
			wantErr: false,
		},
		{
			name:    "Invalid input",
			args:    args{s: string(invalidBlock)},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodeString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DecodeString() = %v, want %v", got, tt.want)
			}
		})
	}
}
