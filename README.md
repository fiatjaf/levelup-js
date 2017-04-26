[![GoDoc](https://godoc.org/github.com/fiatjaf/go-levelup-js?status.png)](http://godoc.org/github.com/fiatjaf/go-levelup-js) [![travis ci badge](https://travis-ci.org/fiatjaf/levelup-js.svg?branch=master)](https://travis-ci.org/fiatjaf/levelup-js)

# levelup bindings for gopherjs

This package implements the [levelup](https://github.com/fiatjaf/levelup) interface, and at the same time wraps https://www.npmjs.com/package/levelup so you can use any of the [supported leveldown backends](https://github.com/Level/levelup/wiki/Modules#storage-back-ends).

## how to use

```go
package main

import (
	"github.com/fiatjaf/go-levelup"
	"github.com/fiatjaf/go-levelup-js"
	"github.com/fiatjaf/levelup/stringlevelup"
)

func main() {
	bdb := levelupjs.NewDatabase("my-beautiful-browser-database", "fruitdown")
	defer bdb.Erase()

	db := stringlevelup.StringDB(bdb)

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

if you don't call `stringlevelup.StringDB` on the object returned by `NewDatabase` you still can use the same methods, only replacing all `string` arguments with `[]byte` (returned values will also be bytes). I've just used the string approach here for readability, but it is slower.

`levelupjs` will try to `require()` both `levelup` and the name of the adapter you gave. If `require` is not available it will try to load these from the global namespace, so `window.levelup` and `window[adapterName]` must be set. If you're running this on a Node.js environment you should be fine with the above instructions, but if you're just using raw HTML you'll have to include these dependencies before, like:

## How to use:

### plain browser <script>

You can just include `levelup` and the adapter of your choice. They will be picked from the global namespace.

```html
<!doctype html>

<script src=https://wzrd.in/standalone/levelup></script>
<script src=https://wzrd.in/standalone/fruitdown></script>
<script src=your-script-that-includes-levelupjs-somehow.js></script>
```

If you're using Webpack, this is probably what you have to do.

### [browserify](http://browserify.org/)

For the JS libraries to be available for being required from the GopherJS environment they must be "globally" requireable, that's not possible if you do the normal browserifying stuff, but it is if you bundle them the following way:

```shell
go get github.com/fiatjaf/go-levelup-js
npm install levelup
npm install fruitdown # or other adapter, you probably want 'memdown' for testing.
browserify your-script.js -r levelup -r memdown -o bundle.js
```

Then include `bundle.js` in the HTML.

### standalone program (backed by Node.js)

```shell
go get github.com/fiatjaf/go-levelup-js
npm install levelup
npm install fruitdown # or other adapter, you probably want 'memdown' for testing.
gopherjs run your-script.go
```
