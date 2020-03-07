package gowasmhelper

import (
	"errors"
	"sync"
	"syscall/js"
)

var (
	errNoJsValueId = errors.New("A JS Value id is required")
)

type fulFullFunc func(event eventData) error

type element struct {
	mux sync.Mutex

	// this A pointer to the js.Value in the Dom
	this *js.Value

	// subname required - The sub name string is used to generate a unique id
	// it must be unique within the array of child nodes.
	subname string

	// The html tag to use
	name htmlTag

	// The id - this will be auto populated always
	id string

	// role is optional
	role string

	// attributes are optional
	attributes map[string]string

	// attributes are optional
	class []string

	// text is the value of the time in display
	text prop

	// used for tags like the image source
	src prop

	// type of the element
	eType prop

	// the event names that the element is registered for
	registeredEvents []jsCallback

	// If the property is true the value of subname is registered as the
	// element id which becomes the dom element id. If this value is false (default)
	// a distinct id is generated from the path.
	useSubNameAsId bool

	// contains the go code for fulfilling the event
	fulFill fulFullFunc

	// contentTarget holds the absolute path of the element
	// for which fulfillment may modify
	contentTarget string

	// children is used to define sub elements. the page api process children in
	// depthwise order when the page is built
	children []*element
}

func (ele *element) addClass(name string) {
	addJSValueClass(ele.this, name)
	return
}

func (ele *element) classExists(name string) bool {
	return existJSValueClass(ele.this, name)
}

func (ele *element) toggleClass(name string) (*js.Value, error) {
	return toggleJSValueClass(ele.this, name)
}

func (ele *element) setInnerHTML(html string) error {
	return pageControl.updateElementOnPage(ele.contentTarget, html)
}

func (ele *element) clearTarget() error {
	return pageControl.clearElementOnPage(ele.contentTarget)
}

func (ele *element) getId() string {
	return ele.this.Get(elementId).String()
}

func (ele *element) createJSValue() error {
	var err error

	ele.this, err = createJSElement(ele.name)
	if err != nil {
		return err
	}

	if ele.id == EMPTYSTRING {
		return errNoJsValueId
	}

	// TODO: Update the js Value properties to match the element
	return nil
}

func (ele *element) appendChild(childElement *element) error {
	mergeElementWithJSValue(childElement)
	ele.this.Call(jsAC, childElement.this)
	return nil
}

func mergeElementWithJSValue(ele *element) {
	// Set the ID
	ele.this.Set(elementId, ele.id)

	// Set the inner html
	ele.this.Set(innerHTML, ele.text)

	// Add the classes
	addJSValueClass(ele.this, ele.class...)

	// Update the source value if provided
	if ele.src != "" {
		ele.this.Set(src, ele.src)
	}

	// Update the type value if provided
	if ele.eType != "" {
		ele.this.Set(eType, ele.eType)
	}

	// Register any call backs provided
	for _, cbItem := range ele.registeredEvents {
		addJSEventListener(ele.this, cbItem.method, cbItem.callback)
	}
}
