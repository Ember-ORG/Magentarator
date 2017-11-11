package main

import (
	"honnef.co/go/js/dom"
)

func main() {
	d := dom.GetWindow().Document()
	foo := d.GetElementByID("create").(*dom.HTMLButtonElement)
	foo.AddEventListener("click", false, func(event dom.Event) {
		println("Test")
	})
}
