package wasmui

//
// import (
// 	"errors"
// 	"path"
// 	"strconv"
// 	"strings"
// 	"sync"
// 	"syscall/js"
// )
//
// var pageControl *page
// var pageBody *element
//
// var panicOnErr = func(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }
//
// func init() {
// 	pageControl = &page{
// 		id:      "document",
// 		this:    js.Global().Get("document"),
// 		pageMap: make(map[string]*element),
// 	}
//
// 	jsVal, err := getJSValueById("body")
// 	if err != nil {
// 		log("%s", err.Error())
// 	}
// 	panicOnErr(err)
//
// 	_, err = pageControl.addElement("html", &element{
// 		name:    body,
// 		subname: body,
// 		this:    jsVal,
// 	})
// 	if err != nil {
// 		log("%s", err.Error())
// 	}
// }
//
// var (
// 	errSubName = errors.New("Subname property is required")
// )
//
// type page struct {
// 	mux           sync.Mutex
// 	id            string
// 	this          js.Value
// 	pageMap       map[string]*element
// 	usr           *user
// 	loginRequired js.Func
// 	fulFillMap    map[string]func(ed eventData) error
// }
//
// func (p *page) getElementsByPrefix(ePath string) ([]*element, error) {
// 	eles := []*element{}
// 	for k, v := range p.pageMap {
// 		if strings.HasPrefix(k, ePath) {
// 			eles = append(eles, v)
// 		}
// 	}
// 	return eles, nil
// }
//
// func (p *page) getElementByPath(ePath string) (*element, error) {
// 	if val, ok := p.pageMap[ePath]; ok {
// 		if val.this != nil {
// 			if val.this.Type() != js.TypeNull {
// 				return val, nil
// 			}
// 		}
// 	}
// 	return nil, errNotFound
// }
//
// func (p *page) clearElementOnPage(ePath string) error {
// 	return p.updateElementOnPage(ePath, EMPTYSTRING)
// }
//
// func (p *page) updateElementOnPage(ePath string, html string) error {
// 	if _, ok := p.pageMap[ePath]; ok {
// 		p.mux.Lock()
// 		p.pageMap[ePath].this.Set(innerHTML, html)
// 		p.mux.Unlock()
// 		return nil
// 	}
// 	return errNotFound
// }
//
// func (p *page) addToPageMap(ePath string, ele *element) {
// 	p.mux.Lock()
// 	p.pageMap[ePath] = ele
// 	p.mux.Unlock()
// }
//
// func (p *page) addElement(parentPath string, ele *element) (string, error) {
// 	var err error
// 	if parentPath == "html" {
// 		p.addToPageMap(path.Join(parentPath, body), ele)
// 		return "", nil
// 	}
//
// 	// Get the parent
// 	parent, err := p.getElementByPath(parentPath)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// Sub name must be provided to create a unique id for the element
// 	if ele.subname == "" {
// 		return "", errSubName
// 	}
//
// 	if !ele.useSubNameAsId {
// 		// Set the new id of the element
// 		ele.id = getElementHashId(path.Join(parentPath, ele.subname, ele.name))
// 	} else {
// 		ele.id = ele.subname
// 	}
//
// 	// Create the js element to be appended
// 	if err := ele.createJSValue(); err != nil {
// 		return "", err
// 	}
//
// 	// Get the reference path of the id
// 	elePath := path.Join(parentPath, ele.id)
//
// 	//log("%+v",elePath)
// 	// Add it to the dom first
// 	err = parent.appendChild(ele)
// 	if err != nil {
// 		return "", err
// 	}
//
// 	// Add it to the internal map
// 	p.addToPageMap(elePath, ele)
//
// 	return elePath, nil
// }
//
// // pushes items to display and returns an empty container
// func pushLayout(eles []*element, parent string) []*element {
// 	for pos, item := range eles {
// 		if !item.useSubNameAsId {
// 			item.subname = strconv.Itoa(pos)
// 		}
//
// 		path, err := pageControl.addElement(parent, item)
// 		if err != nil {
// 			log("%s - %s", err.Error(), parent)
// 			for k := range pageControl.pageMap {
// 				log("%s", k)
// 			}
// 			return nil
// 		}
//
// 		// Add to allow layouts to define children
// 		// this is done at build time. Considering a ring implementation in the
// 		// future
// 		if item.children != nil {
// 			pushLayout(item.children, path)
// 		}
// 	}
// 	return []*element{}
// }
//
// func newBaseElement(parentId string, name htmlTag, class []string, basePath string) (*element, string, error) {
//
// 	parentPath := path.Join(basePath, parentId)
// 	ele, err := pageControl.getElementByPath(parentPath)
// 	if err != nil {
// 		return nil, "", err
// 	}
// 	ele.children = []*element{
// 		&element{
// 			name:  name,
// 			class: class,
// 		},
// 	}
// 	return ele, parentPath, nil
// }
//
// func newClickCallback() []jsCallback {
// 	return []jsCallback{
// 		jsCallback{
// 			method:   jsClick,
// 			callback: pageEvents.dispatch,
// 		},
// 	}
// }
//
// func actionFulfillment(ed eventData) error {
// 	if _, ok := pageControl.fulFillMap[ed.path]; !ok {
// 		log("Default Fulfillment:%+v", ed)
// 	}
// 	return nil
// }
