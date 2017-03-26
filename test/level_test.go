package main

import (
	level "github.com/fiatjaf/go-levelup-js"
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/console"
)

func main() {
	db := level.NewDatabase("", js.Global.Get("memdown"))
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
	db.Batch([]level.Operation{
		level.OpPut("key2", "w"),
		level.OpPut("key3", "z"),
		level.OpDel("key1"),
		level.OpPut("key1", "t"),
		level.OpPut("key4", "m"),
		level.OpPut("key5", "n"),
		level.OpDel("key3"),
	})
	res, _ = db.Get("key1")
	console.Log("res at key1: ", res)
	res, _ = db.Get("key2")
	console.Log("res at key2: ", res)
	res, _ = db.Get("key3")
	console.Log("res at key3: ", res)

	console.Log("reading all")
	iter := db.ReadRange(level.RangeOpts{})
	for iter.Next() {
		console.Log("row: ", iter.Key(), " ", iter.Value())
	}
	console.Log("iter error: ", iter.Error())
}
