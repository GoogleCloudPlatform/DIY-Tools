package domhelper

import (
	"fmt"
	"hash/fnv"
	"path"

	"strings"
	"syscall/js"
)

// Not sure if I need the below.
// type wasmError struct {
// 	msg          string
// 	replacements []string
// }
//
// func (e *wasmError) Error() string {
// 	return fmt.Sprintf(msg, replacements...)
// }

var (
	document js.Value

	// properties
	//elementId     prop = "id"
	//role          prop = "role"

	// these seem like methods because don't want to copy these values
	parentElement prop = "parentElement"
	childNodes    prop = "childNodes"
	innerHTML     prop = "innerHTML"

	classList prop = "classList"
	src       prop = "src"
	eType     prop = "type"
	domText   prop = "text"
	length    prop = "length"

	//tags
	div         htmlTag = "div"
	body        htmlTag = "body"
	script      htmlTag = "script"
	header      htmlTag = "header"
	footer      htmlTag = "footer"
	button      htmlTag = "button"
	input       htmlTag = "input"
	mainContent htmlTag = "main"
	nav         htmlTag = "nav"
	img         htmlTag = "img"
	span        htmlTag = "span"

	//Methods
	jsCE       jsMethod = "createElement"
	jsAC       jsMethod = "appendChild"
	jsGEByID   jsMethod = "getElementById"
	jsGEByName jsMethod = "getElementByName"
	jsCN       jsMethod = "childNodes"
	jsSA       jsMethod = "setAttribute"
	jsAdd      jsMethod = "add"
	jsAEL      jsMethod = "addEventListener"
	jsToggle   jsMethod = "toggle"
	jsContains jsMethod = "contains"

	//events
	jsClick   jsEvent = "click"
	jsOnInput jsEvent = "input"
	jsKeyUp   jsEvent = "keyup"
)

type htmlTag = string
type jsMethod = string
type jsClass = string

type prop = string
type jsEvent = string

// GetChildNodeCount return the count of children dom node.
func GetChildNodeCount(jsVal *js.Value) int {
	return jsVal.Get(childNodes).Get(length).Int()
}

// GetJSValueId gets the ID of the DOM element. It return empty string if the element has not ID.
func GetJSValueId(ele *js.Value) string {
	return ele.Get(elementId).String()
}

// GetJSValueByName gets a list of elements that match the requested name
func GetJSValueByName(name htmlTag) (*js.Value, error) {
	return jsMethodCall(nil, jsGEByName, NOTFOUND, name)
}

// func getJSValueByHashedPath(ePath string) (*js.Value, error) {
// 	id := getElementHashId(ePath)
// 	log("ePath:%s id:%s", ePath, id)
// 	return jsMethodCall(nil, jsGEByID, NOTFOUND, id)
// }

// GetJSValueById get the jsvalue for a given ID. An Error will be returned if the item is not found.
func GetJSValueById(id string) (*js.Value, error) {
	return jsMethodCall(nil, jsGEByID, NOTFOUND, id)
}

// GetParentJSValue get the parent jsValue for a given jsValue. An error is returned if no parent is found.
func GetParentJSValue(e *js.Value) (*js.Value, error) {
	p := e.Get(parentElement)
	if p.Type() == js.TypeNull {
		return nil, errNotFound
	}
	return &p, nil
}

// CreateJSElement create a jsValue with the type given in the name. Example (div, span, etc)
func CreateJSElement(name prop) (*js.Value, error) {
	return jsMethodCall(nil, jsCE, "", name)
}

// SetJSElementId set the dom elementid of the jsValue
func SetJSElementId(data string, jsVal *js.Value) (*js.Value, error) {
	return jsMethodCall(jsVal, jsSA, "", data)
}

// AddJSValueClass add the class to the class list
func AddJSValueClass(jsVal *js.Value, class ...string) {
	for _, c := range class {
		clist := jsVal.Get(classList)
		jsMethodCall(&clist, jsAdd, "", c)
	}
}

// ExistJSValueClass checks if a class exists in the existing element
func ExistJSValueClass(jsVal *js.Value, class string) bool {
	cList := jsVal.Get(classList)

	j, err := jsMethodCall(&cList, jsContains, "", class)
	if err != nil {
		return false
	}
	if j.Type() == js.TypeBoolean {
		return j.Bool()
	}
	return false
}

// ToggleJSValueClass toggle the class of an existing element
func ToggleJSValueClass(jsVal *js.Value, class string) (*js.Value, error) {
	clist := jsVal.Get(classList)
	return jsMethodCall(&clist, jsToggle, "", class)
}

// AddJSEventListener adds an event listener to a dom element
func AddJSEventListener(jsVal *js.Value, eventName jsEvent, cb js.Func) {
	_, err := jsMethodCall(jsVal, jsAEL, "", eventName, cb)
	if err != nil {
		panic(err)
	}
}

func jsMethodCall(jsVal *js.Value, method string, errtxt string, params ...interface{}) (*js.Value, error) {
	if jsVal == nil {
		jsVal = &pageControl.this
	}
	val := jsVal.Call(method, params...)
	if val.Type() == js.TypeNull {
		return nil, fmt.Errorf(errtxt, params...)
	}
	return &val, nil
}

func getElementHashId(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	return fmt.Sprint(hash.Sum32())
}

// getTargetPath get the containing element ids as a route.
// The idea is to be able to use this route take action on button click
func getTargetPath(ele *element) (string, error) {
	ids := []string{}
	reverse := func(a []string) {
		for i := len(a)/2 - 1; i >= 0; i-- {
			opp := len(a) - 1 - i
			a[i], a[opp] = a[opp], a[i]
		}
	}

	// Try the get the element in the dom by the existing id
	var jsElem *js.Value
	var err error
	if jsElem, err = getJSValueById(ele.id); err == nil {
		ids = append(ids, jsElem.Get(elementId).String())
		for {
			if jsElem, err = getParentJSValue(jsElem); err == errNotFound {
				break
			}
			ids = append(ids, jsElem.Get(elementId).String())
		}
	} else {
		return "", fmt.Errorf(NOTFOUND, ele)
	}

	// reverse ids so we get the paths in the right order
	reverse(ids)

	// join them and return
	thePath := path.Clean(strings.Join(ids, "/"))
	return thePath, nil
}
