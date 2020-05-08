package benchmarks

import (
	"bytes"
	"io/ioutil"
	"testing"

	own "github.com/akshaybharambe14/go-jsonc"
	jsonc "muzzammil.xyz/jsonc"
)

var jsonBytesSmall, _ = ioutil.ReadFile("./testdata/small.jsonc")
var jsonBytesBig, _ = ioutil.ReadFile("./testdata/big.jsonc")

func BenchmarkOwnSmallJSONBytes(b *testing.B) {
	cnt := len(jsonBytesSmall)
	b.SetBytes(int64(cnt))

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		ip := make([]byte, cnt)
		copy(ip, jsonBytesSmall)

		b.StartTimer()

		_, _ = own.DecodeBytes(ip)
	}
}

func BenchmarkOwnSmallJSONBytesReader(b *testing.B) {
	b.SetBytes(int64(len(jsonBytesSmall)))

	for i := 0; i < b.N; i++ {
		_, _ = ioutil.ReadAll(own.NewDecoder(bytes.NewBuffer(jsonBytesSmall)))
	}
}

func BenchmarkJSONCSmallJSONBytes(b *testing.B) {
	b.SetBytes(int64(len(jsonBytesSmall)))

	for i := 0; i < b.N; i++ {
		jsonc.ToJSON(jsonBytesSmall)
	}
}

func BenchmarkOwnBigJSONBytes(b *testing.B) {
	cnt := len(jsonBytesBig)
	b.SetBytes(int64(cnt))

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		ip := make([]byte, len(jsonBytesBig))
		copy(ip, jsonBytesBig)

		b.StartTimer()

		_, _ = own.DecodeBytes(ip)
	}
}

func BenchmarkOwnBigJSONBytesReader(b *testing.B) {
	b.SetBytes(int64(len(jsonBytesBig)))

	for i := 0; i < b.N; i++ {
		_, _ = ioutil.ReadAll(own.NewDecoder(bytes.NewBuffer(jsonBytesBig)))
	}
}

func BenchmarkJSONCBigJSONBytes(b *testing.B) {
	b.SetBytes(int64(len(jsonBytesBig)))

	for i := 0; i < b.N; i++ {
		_ = jsonc.ToJSON(jsonBytesBig)
	}
}
