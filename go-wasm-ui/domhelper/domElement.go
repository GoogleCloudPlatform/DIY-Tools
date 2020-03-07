package domHelper

import(
//  "github.com/goog-lukemc/gowasmui"
)

type Element stuct{
  // ID is the id of the HTML element. In this implementation it is required that every
  // element have an ID and they IDs are unique within the page.
  ID string

  // typ is usedful to set at create type of the elements
  typ htmlTag

  // text is the initial inner text of the element. This property is only used to setAttribute
  // the initial value
  Text string
}

// HTMLTag is create to contol which HTML tag can be used in the solution. The tags below are tested.
type HTMLTag struct{
  Body string
  Div string
  Script string
  Header string
  Footer string
  Input string
  Img string
}
