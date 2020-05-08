# JSON with comments for GO

- Decodes a "commented json" to "json". Provided, the input must be a valid jsonc document.
- Supports io.Reader
- With this, we can use commented json files as configuration for go applications.

Inspired by [muhammadmuzzammil1998](https://github.com/muhammadmuzzammil1998/jsonc)

```jsonc
{
	/*
        some block comment
  */
	"string": "foo", // a string
	"bool": false, // a boolean
	"number": 42, // a number
	// "object":{
	//     "key":"val"
	// },
	"array": [
		// example of an array
		1,
		2,
		3
	]
}
```

Gets converted to (spaces omitted)

```json
{ "string": "foo", "bool": false, "number": 42, "array": [1, 2, 3] }
```

## Motivation

[jsonc](https://github.com/muhammadmuzzammil1998/jsonc) is great. But this package provides significant performance improvements and simple API to use it with standard library. See [benchmarks](https://link)

## Usage

Get this package

```sh

go get github.com/akshaybharambe14/go-jsonc

```

## Benchmarks

Here is the performance comparison as compared to [JSONC](https://github.com/muhammadmuzzammil1998/jsonc)

go-jsonc avoids allocations to heap to gain the performance.

```text
goos: windows
goarch: amd64
pkg: github.com/akshaybharambe14/go-jsonc/benchmarks
BenchmarkOwnSmallJSONBytes-4              256599              4952 ns/op         353.00 MB/s           0 B/op          0 allocs/op
BenchmarkOwnSmallJSONBytesReader-4        206823              5832 ns/op         299.70 MB/s        6224 B/op          5 allocs/op
BenchmarkJSONCSmallJSONBytes-4            171474              6925 ns/op         252.41 MB/s        1792 B/op          1 allocs/op
BenchmarkOwnBigJSONBytes-4                 33517             35921 ns/op         462.26 MB/s           0 B/op          0 allocs/op
BenchmarkOwnBigJSONBytesReader-4          105244             11292 ns/op        1470.45 MB/s        6224 B/op          5 allocs/op
BenchmarkJSONCBigJSONBytes-4               19599             61422 ns/op         270.34 MB/s       18432 B/op          1 allocs/op
PASS
ok      github.com/akshaybharambe14/go-jsonc/benchmarks 26.250s
```

## Example

see [examples](https://github.com/akshaybharambe14/go-jsonc/tree/master/examples)

## License

`go-jsonc` is open source and available under [MIT License](License.md)
