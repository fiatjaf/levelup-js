# levelup bindings for gopherjs.  [![GoDoc](https://godoc.org/github.com/fiatjaf/go-levelup-js?status.png)](http://godoc.org/github.com/fiatjaf/go-levelup-js) [![travis ci badge](https://travis-ci.org/fiatjaf/levelup-js.svg?branch=master)](https://travis-ci.org/fiatjaf/levelup-js)


## How to use:

```
go get github.com/fiatjaf/go-levelup-js
```


```go
package main

import (
	"github.com/fiatjaf/go-levelup"
	"github.com/fiatjaf/go-levelup-js"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

func main() {
	updb := levelupjs.NewDatabase("dbname", js.Global.Get("fruitdown"))
	db := stringlevelup.StringDB(updb)

	fmt.Println("setting key1 to x")
	db.Put("key1", "x")
	res, _ := db.Get("key1")
	fmt.Println("setting key2 to 2")
	fmt.Println("res at key2: ", res)
	db.Put("key2", "y")
	res, _ = db.Get("key2")
	fmt.Println("res at key2: ", res)
	fmt.Println("deleting key1")
	db.Del("key1")
	res, _ = db.Get("key1")
	fmt.Println("res at key1: ", res)

	fmt.Println("batch")
	db.Batch([]levelup.Operation{
		stringlevelup.Put("key2", "w"),
		stringlevelup.Put("key3", "z"),
		stringlevelup.Del("key1"),
		stringlevelup.Put("key1", "t"),
		stringlevelup.Put("key4", "m"),
		stringlevelup.Put("key5", "n"),
		stringlevelup.Del("key3"),
	})
	res, _ = db.Get("key1")
	fmt.Println("res at key1: ", res)
	res, _ = db.Get("key2")
	fmt.Println("res at key2: ", res)
	res, _ = db.Get("key3")
	fmt.Println("res at key3: ", res)

	fmt.Println("reading all")
	iter := db.ReadRange(nil)
	for ; iter.Valid(); iter.Next() {
		fmt.Println("row: ", iter.Key(), " ", iter.Value())
	}
	fmt.Println("iter error: ", iter.Error())
	iter.Release()
```

For now the library expects to find `levelup` in the global object (`window` in browser).
For all adapters you're using it is also expected that the name you pass to `NewDatabase` corresponds to the global name of the adapter.

```html
<!doctype html>

<script src=https://wzrd.in/standalone/levelup></script>
<script src=https://wzrd.in/standalone/fruitdown></script>
<script src=bundle.js></script>
```
