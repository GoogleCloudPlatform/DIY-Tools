package domhelper

import (
	"fmt"
	"path"
	"runtime"
)

// CLog is a modified version of a standard console log. It provides additional runtime information
// to the console log in the browser output.
func CLog(msg string, objs ...interface{}) {

	// Get useful details from the runtime
	fpcs := make([]uintptr, 1)
	runtime.Callers(2, fpcs)
	caller := runtime.FuncForPC(fpcs[0] - 1)
	fileName, position := caller.FileLine(fpcs[0] - 1)
	funcName := caller.Name()

	objs = append(objs, path.Base(fileName), funcName, position)
	//var s string
	//msg = msg + s + "::%v : %v : %v"

	// Assemble a replaceable string
	msg = msg + "::%v : %v : %v"

	// Print the results to the console.
	// TODO: This cloud be the only use of the fmt package. Explore builtin printf()
	// to reduce WASM size.
	fmt.Printf(msg+"\n", objs...)
}
