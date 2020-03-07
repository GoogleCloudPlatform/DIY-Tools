package gowasmhelper

import (
	//"fmt"
	//"path"
	//"strings"
	"syscall/js"
)

var pageEvents domEvent

type jsCallbackData struct {
	this *js.Value
	args []*js.Value
}

type jsCallback struct {
	method   string
	callback js.Func
}

type domEvent struct {
	dispatch  js.Func
	eventChan chan jsCallbackData
}

type eventData struct {
	path      string
	eventType string
	event     *js.Value
	args      []*js.Value
}

func init() {
	domEventSetup()
}

func (cer *domEvent) eventRouter() {
	for ce := range cer.eventChan {
		// Testing out background events
		log("Type Is %s", ce.this.Type())

		ed := eventData{}

		// Create a shadow element so we can get the element path
		ele := &element{
			id: ce.this.Get(elementId).String(),
		}

		// Get the event type for later used by fulfillment
		var t string
		if len(ce.args) > 0 {
			t = ce.args[0].Get("type").String()
		}

		// get the element path of the calling event
		eventSource, err := getTargetPath(ele)
		if err != nil {
			log("Event worker Error:%s", err.Error())
		}

		// get the element from from the go vdom
		ele, err = pageControl.getElementByPath(eventSource)
		if err != nil {
			log("path:%s error:%s", eventSource, err.Error())
		}

		// If fulfill has been implemented call it with the event data
		if ele.fulFill == nil {
			log("No Fulfillment function for %s", eventSource)
			continue
		}

		// Setting the path of the target.
		ed.path = eventSource

		// setting a type property of conven
		ed.eventType = t

		// setting the event data to the event property
		ed.event = ce.args[0]

		// sending all sent args with the event
		ed.args = ce.args

		// Call the fulfill functions
		go ele.fulFill(ed)
	}
}

func domEventSetup() {
	// Sets up the core dom event router
	pageEvents = domEvent{
		eventChan: make(chan jsCallbackData, 1),
	}

	pageEvents.dispatch = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ce := jsCallbackData{
			this: &this,
		}
		for _, a := range args {
			ce.args = append(ce.args, &a)
		}
		pageEvents.eventChan <- ce
		return nil
	})

	// All event router
	go pageEvents.eventRouter()

}
