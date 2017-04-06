package levelupjs_test

import (
	"testing"

	"github.com/fiatjaf/levelup-js"
	tests "github.com/fiatjaf/levelup/tests"
)

func TestAll(t *testing.T) {
	db := levelupjs.NewDatabase("", "memdown")
	defer db.Erase()

	tests.Test(db, t)
}
