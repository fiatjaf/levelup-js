package levelupjs

import (
	"github.com/fiatjaf/go-levelup"
	"github.com/gopherjs/gopherjs/js"
)

func NewDatabase(dbName string, backend *js.Object) *Level {
	db := js.Global.Get("levelup").Invoke(dbName, map[string]*js.Object{
		"db": backend,
	})
	return &Level{db}
}

type Level struct {
	db *js.Object
}

func (l Level) Put(key, value string) (interface{}, error) {
	rw := make(resultWaiter)
	l.db.Call("put", key, value, rw.Done)
	return rw.Read()
}

func (l Level) Get(key string) (interface{}, error) {
	rw := make(resultWaiter)
	l.db.Call("get", key, rw.Done)
	return rw.Read()
}

func (l Level) Del(key string) (interface{}, error) {
	rw := make(resultWaiter)
	l.db.Call("del", key, rw.Done)
	return rw.Read()
}

func (l Level) Batch(ops []levelup.Operation) (interface{}, error) {
	rw := make(resultWaiter)
	l.db.Call("batch", ops, rw.Done)
	return rw.Read()
}

// ---

func (l Level) ReadRange(opts levelup.RangeOpts) *ReadIterator {
	actualOpts := map[string]interface{}{}
	if opts.GT != "" {
		actualOpts["gt"] = opts.GT
	}
	if opts.GTE != "" {
		actualOpts["gte"] = opts.GTE
	}
	if opts.LT != "" {
		actualOpts["lt"] = opts.LT
	}
	if opts.LTE != "" {
		actualOpts["lte"] = opts.LTE
	}
	if opts.Reverse {
		actualOpts["reverse"] = opts.Reverse
	}
	if opts.Limit > 0 {
		actualOpts["limit"] = opts.Limit
	}

	stream := l.db.Call("createReadStream", actualOpts)
	ri := &ReadIterator{
		stream:  stream,
		current: -1,
		err:     nil,
		open:    true,
		ended:   false,
		emitted: make(chan int),
	}

	ri.stream.Call("on", "data", func(data *js.Object) {
		ri.all = append(ri.all, data)
		ri.emitted <- 1
	})

	ri.stream.Call("on", "error", func(jserr *js.Object) {
		ri.err = JSError{jserr}
		ri.open = false
		ri.ended = true
	})

	ri.stream.Call("on", "close", func(data *js.Object) {
		ri.open = false
	})

	ri.stream.Call("on", "end", func(data *js.Object) {
		ri.ended = true
	})

	return ri
}

type ReadIterator struct {
	stream  *js.Object
	current int
	all     []*js.Object
	err     error
	open    bool
	ended   bool
	emitted chan int
}

func (ri *ReadIterator) Next() bool {
	if ri.open && !ri.ended {
		<-ri.emitted
		ri.current++
		return true
	}
	return false
}

func (ri *ReadIterator) Key() string {
	return ri.all[ri.current].Get("key").String()
}

func (ri *ReadIterator) Value() string {
	return ri.all[ri.current].Get("value").String()
}

func (ri *ReadIterator) Error() error {
	return ri.err
}

// ---

type result struct {
	result *js.Object
	err    *js.Object
}

type resultWaiter chan result

func (rw resultWaiter) Done(err *js.Object, res *js.Object) {
	rw <- result{res, err}
}

func (rw resultWaiter) Read() (*js.Object, error) {
	r := <-rw
	if r.err == nil {
		return r.result, nil
	}
	return r.result, JSError{r.err}
}

type JSError struct {
	*js.Object
}

func (jserr JSError) Error() string {
	if jserr.Object == js.Undefined {
		return "<nil>"
	} else {
		return "JavaScript error: " + jserr.Get("message").String()
	}
}
