package domhelper

//  "github.com/goog-lukemc/gowasmui"
import (
	"fmt"
	"hash/fnv"
	"syscall/js"
)

// Element
type Element struct {
	// ID is the id of the HTML element. In this implementation it is required that every
	// element have an ID and they IDs are unique within the page.
	ID string

	// typ is usedful to set at create type of the elements
	typ HTMLTagNames

	// Text is the initial inner text of the element. This property is only used to setAttribute
	// the initial value
	Text string

	//jsValue is the underlying dom item
	jsValue *js.Value
}

// Create adsad
func (ele *Element) Create(domPath string, typ string, initialText string) error {
	jv, err := jsMethodCall(ele.jsValue, JSMethod.CreateElement, typ)
	ele.jsValue = jv
	return err
}

// SetAttribute das
func (ele *Element) SetAttribute(name string, value string) error {
	if ele.jsValue == nil {
		return fmt.Errorf("noExistingJSValue: existing js value requried to call this method")
	}
	_, err := jsMethodCall(ele.jsValue, JSMethod.SetAttribute, name, value)
	return err
}

// AddEventListener adds an event listener to a dom element
func (ele *Element) AddEventListener(eventName string, cbs ...js.Func) error {
	for _, cb := range cbs {
		_, err := jsMethodCall(ele.jsValue, JSMethod.AddEventListner, eventName, cb)
		if err != nil {
			return err
		}
	}
	return nil
}

func jsMethodCall(jsVal *js.Value, method string, params ...interface{}) (*js.Value, error) {
	if jsVal == nil {
		g := js.Global()
		jsVal = &g
	}
	val := jsVal.Call(method, params...)
	if val.Type() == js.TypeNull {
		return nil, fmt.Errorf("unexpectedNull: a null js value is not expected")
	}
	return &val, nil
}

// newElementIdFromDOMPath dsa
func newElementIDFromDOMPath(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return fmt.Sprint(hash.Sum32())
}

// HTMLTagNames is create to contol which HTML tag can be used in the solution. The tags below are tested.
type HTMLTagNames struct {
	Body   string
	Div    string
	Script string
	Header string
	Footer string
	Input  string
	Img    string
}

// JSMethodNames contrains the correct string property name for tested javascript methods
type JSMethodNames struct {
	CreateElement   string
	AppendChild     string
	GetElementByID  string
	SetAttribute    string
	AddEventListner string
	Toggle          string
	Contains        string
}

// JSEventNames represent the event name tested with the js event handler in this package.
type JSEventNames struct {
	Click   string
	OnInput string
	Keyup   string
}

var (
	// HTMLTag is export from this package as a global variable and represents a property for each support HTMLTagName
	HTMLTag HTMLTagNames

	// JSMethod is export from this package as a global variable and represents a property for each support Javescript Method
	JSMethod JSMethodNames

	// JSEvent is export from this package as a global variable and represents a property for each support Javescript Event
	JSEvent JSEventNames
)

func init() {
	HTMLTag = HTMLTagNames{
		Body:   "body",
		Div:    "div",
		Script: "string",
		Header: "header",
		Footer: "footer",
		Input:  "input",
		Img:    "img",
	}

	JSMethod = JSMethodNames{
		CreateElement:   "createElement",
		AppendChild:     "appendChild",
		GetElementByID:  "getElementByID",
		SetAttribute:    "SetAttribute",
		AddEventListner: "addEventListner",
		Toggle:          "toggle",
		Contains:        "contains",
	}

	JSEvent = JSEventNames{
		Click:   "click",
		OnInput: "onInput",
		Keyup:   "keyUp",
	}

}
