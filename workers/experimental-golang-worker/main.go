package main

import (
	"syscall/js"
)

var (
	fetch   = js.Global().Get("fetch")
	promise = js.Global().Get("Promise")
)

func main() {
	js.Global().Set("goGetter", getter)

	select {}
}

var getter = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	const exampleURL = `https://api.github.com/repos/javascript-tutorial/en.javascript.info/commits`

	ch := await(fetch.Invoke(exampleURL))
	p, set := newPromise()
	go func() {
		results := <-ch
		rsp := results[0]

		results = <-await(rsp.Call("json"))
		set(results[0])
	}()
	return p
})

func newPromise() (p js.Value, set func(js.Value)) {
	ch := make(chan js.Value)
	resolver := make(chan js.Value, 1)
	go func() {
		result := <-ch
		resolve := <-resolver
		resolve.Invoke(result)
	}()
	p = promise.New(js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolver <- args[0]
		return nil
	}))
	set = func(v js.Value) {
		ch <- v
	}
	return
}

func await(awaitable js.Value) chan []js.Value {
	ch := make(chan []js.Value)
	cb := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ch <- args
		return nil
	})
	awaitable.Call("then", cb)
	return ch
}
