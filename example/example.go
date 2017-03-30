package main

import (
	"github.com/fiatjaf/levelup-js"
	examples "github.com/fiatjaf/levelup/examples"
	"github.com/gopherjs/gopherjs/js"
)

func main() {
	db := levelupjs.NewDatabase("", js.Global.Get("memdown"))
	defer db.Erase()

	examples.Example(db)
}
