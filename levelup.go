package levelupjs

import (
	"github.com/fiatjaf/levelup"
	"github.com/gopherjs/gopherjs/js"
)

func NewDatabase(dbName string, backend *js.Object) *Level {
	db := js.Global.Get("levelup").Invoke(
		dbName,
		map[string]*js.Object{
			"db": backend,
		},
	)
	return &Level{db}
}

type Level struct {
	db *js.Object
}

func (l Level) Put(key, value string) error {
	rw := make(resultWaiter)
	l.db.Call("put", key, value, rw.Done)
	_, jserr := rw.Read()
	return jserr.ProbablyNil()
}

func (l Level) Get(key string) (string, error) {
	rw := make(resultWaiter)
	l.db.Call("get", key, rw.Done)
	data, jserr := rw.Read()
	if !jserr.IsNil() {
		if jserr.Name() == "NotFoundError" {
			return "", levelup.NotFound
		}

		return "", jserr.ProbablyNil()
	}
	datastring := data.String()
	if datastring == "" {
		return "", levelup.NotFound
	}
	return datastring, nil
}

func (l Level) Del(key string) error {
	rw := make(resultWaiter)
	l.db.Call("del", key, rw.Done)
	_, jserr := rw.Read()
	return jserr.ProbablyNil()
}

func (l Level) Batch(ops []levelup.Operation) error {
	rw := make(resultWaiter)
	l.db.Call("batch", ops, rw.Done)
	_, jserr := rw.Read()
	return jserr.ProbablyNil()
}

// ---

func (l Level) ReadRange(opts *levelup.RangeOpts) levelup.ReadIterator {
	if opts == nil {
		opts = &levelup.RangeOpts{}
	}
	opts.FillDefaults()

	optsMap := map[string]interface{}{}
	optsMap["gte"] = opts.Start
	optsMap["lt"] = opts.End
	optsMap["reverse"] = opts.Reverse
	optsMap["limit"] = opts.Limit

	// console.Log("reading range ", optsMap)
	stream := l.db.Call("createReadStream", optsMap)
	ri := &ReadIterator{
		cursor: 0,
		err:    JSError{nil},
		open:   true,
		event:  make(chan bool),
	}

	stream.Call("on", "data", func(data *js.Object) {
		// console.Log(" - data! ", data)
		ri.all = append(ri.all, data)
		go func() { ri.event <- true }()
	})

	stream.Call("on", "error", func(jserr *js.Object) {
		// console.Log(" - error! ", jserr)
		ri.err = JSError{jserr}
		ri.open = false
		go func() { ri.event <- true }()
	})

	stream.Call("on", "close", func(data *js.Object) {
		// console.Log(" - close!")
		ri.open = false
		go func() { ri.event <- true }()
	})

	stream.Call("on", "end", func(data *js.Object) {
		// console.Log(" - end!")
		ri.open = false
		go func() { ri.event <- true }()
	})

	return ri
}

type ReadIterator struct {
	cursor int
	all    []*js.Object
	err    JSError
	open   bool
	event  chan bool
}

func (ri *ReadIterator) Valid() bool {
	// setup a timeout because we don't want to wait forever for this
	// it is better to fail altogether
	// timeout := make(chan bool, 1)
	// go func() {
	// 	time.Sleep(time.Second * 25)
	// 	jserr := js.Global.Get("Error").New("[go] timed out while waiting for emitted events.")
	// 	ri.err = JSError{jserr}
	// 	ri.open = false
	// 	timeout <- true
	// }()

	// console.Log("is it valid? ")
	for {
		// we're waiting until the desired position is reached.
		if !ri.open {
			// console.Log("invalid!")
			break
		}

		// if we have already received the value for the current position
		if len(ri.all) > ri.cursor {
			// console.Log("valid!")
			return true
		}

		// this timeout thing must wait until https://github.com/gopherjs/gopherjs/issues/478 is solved.
		// select {
		// case <-ri.event:
		// 	continue
		// case <-timeout:
		// 	break
		// }

		<-ri.event
	}

	// close(ri.event)
	return false
}

func (ri *ReadIterator) Next() {
	if ri.open {
		ri.cursor++
	}
}

func (ri *ReadIterator) Key() string {
	return ri.all[ri.cursor].Get("key").String()
}

func (ri *ReadIterator) Value() string {
	return ri.all[ri.cursor].Get("value").String()
}

func (ri *ReadIterator) Error() error { return ri.err.ProbablyNil() }

func (ri *ReadIterator) Release() {}

// ---

type result struct {
	result *js.Object
	err    *js.Object
}

type resultWaiter chan result

func (rw resultWaiter) Done(err *js.Object, res *js.Object) {
	rw <- result{res, err}
}

func (rw resultWaiter) Read() (*js.Object, JSError) {
	r := <-rw
	return r.result, JSError{r.err}
}

type JSError struct {
	*js.Object
}

func (jserr JSError) ProbablyNil() error {
	if jserr.IsNil() {
		return nil
	}
	return jserr
}
func (jserr JSError) IsNil() bool {
	if jserr.Object == nil {
		return true
	}
	if jserr.Object == js.Undefined {
		return true
	}
	return false
}
func (jserr JSError) Name() string    { return jserr.Get("name").String() }
func (jserr JSError) Message() string { return jserr.Get("message").String() }
func (jserr JSError) Error() string {
	if jserr.Object == js.Undefined {
		return "<nil>"
	} else {
		return "javascript error: " + jserr.Get("message").String()
	}
}
