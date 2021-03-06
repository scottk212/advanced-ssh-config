// +build js

package syscall

import (
	"bytes"
	"unsafe"

	"github.com/gopherjs/gopherjs/js"
)

var warningPrinted = false
var lineBuffer []byte

func init() {
	js.Global.Set("$flushConsole", js.InternalObject(func() {
		if len(lineBuffer) != 0 {
			js.Global.Get("console").Call("log", string(lineBuffer))
			lineBuffer = nil
		}
	}))
}

func printWarning() {
	if !warningPrinted {
		js.Global.Get("console").Call("error", "warning: system calls not available, see https://github.com/gopherjs/gopherjs/blob/master/doc/syscalls.md")
	}
	warningPrinted = true
}

func printToConsole(b []byte) {
	goPrintToConsole := js.Global.Get("goPrintToConsole")
	if goPrintToConsole != js.Undefined {
		goPrintToConsole.Invoke(js.InternalObject(b))
		return
	}

	lineBuffer = append(lineBuffer, b...)
	for {
		i := bytes.IndexByte(lineBuffer, '\n')
		if i == -1 {
			break
		}
		js.Global.Get("console").Call("log", string(lineBuffer[:i])) // don't use println, since it does not externalize multibyte characters
		lineBuffer = lineBuffer[i+1:]
	}
}

func use(p unsafe.Pointer) {
	// no-op
}
