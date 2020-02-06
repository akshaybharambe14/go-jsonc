package jsonc

import (
	"bytes"
	"testing"
)

func ts(b []byte) *Decoder { return &Decoder{r: bytes.NewBuffer(b)} }

var (
	validSingle   = []byte(`{"foo": // this is a single line comment\n"bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com" }`)
	invalidSingle = []byte(`{"foo": // this is a single line comment "bar foo", "true": false, "number": 42, "object": { "test": "done" }, "array" : [1, 2, 3], "url" : "https://github.com" }`)

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
