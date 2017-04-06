package main

import (
	"github.com/fiatjaf/levelup-js"
	examples "github.com/fiatjaf/levelup/examples"
)

func main() {
	db := levelupjs.NewDatabase("", "memdown")
	defer db.Erase()

	examples.Example(db)
}
