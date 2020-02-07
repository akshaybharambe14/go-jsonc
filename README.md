# JSON with comments for GO

- Decodes a "commented json" to "json". Provided, the input must be a valid jsonc document.
- Supports io.Reader

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

## Usage

Get this package

```sh

go get github.com/akshaybharambe14/go-jsonc

```

## Example

see [examples](https://github.com/akshaybharambe14/go-jsonc/examples)

## License

`go-jsonc` is available under [MIT License](License.md)
