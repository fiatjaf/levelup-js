package main

import (
	"github.com/fiatjaf/go-levelup"
	"github.com/fiatjaf/go-levelup-js"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

func main() {
	db := levelupjs.NewDatabase("", js.Global.Get("memdown"))
	console.Log("setting key1 to x")
	res, _ := db.Put("key1", "x")
	res, _ = db.Get("key1")
	console.Log("setting key2 to 2")
	console.Log("res at key2: ", res)
	res, _ = db.Put("key2", "y")
	res, _ = db.Get("key2")
	console.Log("res at key2: ", res)
	console.Log("deleting key1")
	res, _ = db.Del("key1")
	res, _ = db.Get("key1")
	console.Log("res at key1: ", res)

	console.Log("batch")
	db.Batch([]levelup.Operation{
		levelup.OpPut("key2", "w"),
		levelup.OpPut("key3", "z"),
		levelup.OpDel("key1"),
		levelup.OpPut("key1", "t"),
		levelup.OpPut("key4", "m"),
		levelup.OpPut("key5", "n"),
		levelup.OpDel("key3"),
	})
	res, _ = db.Get("key1")
	console.Log("res at key1: ", res)
	res, _ = db.Get("key2")
	console.Log("res at key2: ", res)
	res, _ = db.Get("key3")
	console.Log("res at key3: ", res)

	console.Log("reading all")
	iter := db.ReadRange(levelup.RangeOpts{})
	for iter.Next() {
		console.Log("row: ", iter.Key(), " ", iter.Value())
	}
	console.Log("iter error: ", iter.Error())
}
