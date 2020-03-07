package wasmui

import (
	"fmt"

	"github.com/GoogleCloudPlatform/DIY-Tools/go-wasm-ui/domhelper"
)

type DOMUI struct {
	Head   domhelper.Element
	Body   domhelper.Element
	Foot   domhelper.Element
	Script domhelper.Element
}

// NewApp Creates a new grid on a page
func NewApp() (*DOMUI, error) {
	return &DOMUI{}, nil
}

func SayHi() {
	fmt.Println("Hi From wasmui")
}
