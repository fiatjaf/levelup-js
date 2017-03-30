package levelupjs_test

import (
	"testing"

	"github.com/fiatjaf/levelup-js"
	tests "github.com/fiatjaf/levelup/tests"
	"github.com/gopherjs/gopherjs/js"
)

func TestAll(t *testing.T) {
	db := levelupjs.NewDatabase("", js.Global.Get("memdown"))
	defer db.Erase()

	tests.Test(db, t)
}
