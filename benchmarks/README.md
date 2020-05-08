# Benchmarks

Here is the performance comparison as compared to [jsonc](https://github.com/muhammadmuzzammil1998/jsonc)

go-jsonc avoids allocations to heap to gain the performance.

```
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
