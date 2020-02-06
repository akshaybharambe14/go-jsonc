package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/akshaybharambe14/go-jsonc"
)

func main() {
	// create new decoder
	d := jsonc.NewDecoder(bytes.NewBuffer([]byte(`{
		/*
			some block comment
		*/
		"string": "foo", // a string
		"bool": true, // a boolean
		"number": 42, // a number
		// "object":{
		//     "key":"val"
		// },
		"array": [ // example of an array
			1,
			2,
			3
		]
	}`)))

	// read json
	res, err := ioutil.ReadAll(d)
	if err != nil {
		fmt.Println("error decoding commented json: ", err)
		return
	}

	// Unmarshal to a struct

	type Test struct {
		Str string `json:"string"`
		Bln bool   `json:"bool"`
		Num int    `json:"number"`
		Arr []int  `json:"array"`
	}

	t := Test{}
	if err = json.Unmarshal(res, &t); err != nil {
		fmt.Println("error while json Unmarshal: ", err)
		return
	}

	fmt.Printf("%+v\n", t)
}
