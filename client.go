package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
	"honnef.co/go/js/dom"
)

var conn net.Conn
var aorg = "a"

func main() {
	js.Global.Set("client", map[string]interface{}{
		"start":    start,
		"artclick": artclick,
		"genclick": genclick,
		"add":      add,
		"remove":   remove,
		"onCreate": onCreate,
	})
}

func start() {
	go func() {
		//Defining dom
		d := dom.GetWindow().Document()

		//Auto-selecting artist button
		d.GetElementByID("option-1").(*dom.HTMLInputElement).Click()

		// Establishing server connection
		conn, _ = websocket.Dial("ws://localhost:9001/ws") // Blocks until connection is established.	//On successful connection

		go func() {
			for true {
				buf := make([]byte, 1024)
				n, err := conn.Read(buf) // Blocks until a WebSocket frame is received.
				if string(buf[:n]) == "doneDownloading" {
					d.GetElementByID("progress").SetInnerHTML("Loading Training Data...")
				} else if string(buf[:n]) == "training" {
					d.GetElementByID("progress").SetInnerHTML("Training...")
				} else if string((buf[:n])[0]) == "p" {
					p := strings.Replace(string((buf[:n])), "p", "", -1)
					percent, _ := strconv.ParseFloat(p, 64)
					percent = percent * 100
					js.Global.Call("progress", percent)
					d.GetElementByID("generate").SetAttribute("style", "display: block;")
					d.GetElementByID("p2").SetAttribute("style", "display: none;")
					d.GetElementByID("p1").SetAttribute("style", "display: block;")
				} else if string(buf[:n]) == "generating" {
					d.GetElementByID("progress").SetInnerHTML("Generating music...")
					d.GetElementByID("p1").SetAttribute("style", "display: none;")
					d.GetElementByID("p2").SetAttribute("style", "display: block;")
					d.GetElementByID("generate").SetAttribute("style", "display: none;")
				}

				if err != nil {
					fmt.Println(err.Error())
				}
				time.Sleep(time.Second * 5)
			}
		}()

		//When stop the madness is clicked
		stop := d.GetElementByID("generate").(*dom.HTMLButtonElement)
		stop.AddEventListener("click", false, func(event dom.Event) {
			conn.Write([]byte("fgenerate"))
			d.GetElementByID("progress").SetInnerHTML("Generating Music...")
		})
	}()
}

func artclick() {
	//When artist button is clicked
	d := dom.GetWindow().Document()
	lbls := d.GetElementsByClassName("artistlbl")
	for i, _ := range lbls {
		d.GetElementsByClassName("artistlbl")[i].SetInnerHTML("Artist...")
	}
	aorg = "a"
	fmt.Println(aorg)
}

func genclick() {
	//When artist button is clicked
	d := dom.GetWindow().Document()
	lbls := d.GetElementsByClassName("artistlbl")
	for i, _ := range lbls {
		d.GetElementsByClassName("artistlbl")[i].SetInnerHTML("Genre...")
	}
	aorg = "g"
	fmt.Println(aorg)
}

func add() {
	d := dom.GetWindow().Document()
	d.GetElementByID("remove").(*dom.HTMLButtonElement).Disabled = false
	form := d.CreateElement("form")
	form.SetAttribute("action", "#")
	div := d.CreateElement("div").(*dom.HTMLDivElement)
	div.SetAttribute("class", "mdl-textfield mdl-js-textfield mdl-textfield--floating-label inpt")
	input := d.CreateElement("input")
	input.SetAttribute("class", "mdl-textfield__input artist")
	input.SetAttribute("type", "text")
	input.SetAttribute("onkeydown", "if (event.keyCode == 13){client.onCreate()}")
	label := d.CreateElement("label")
	label.SetAttribute("class", "mdl-textfield__label artistlbl")
	label.SetAttribute("for", "artist")
	aorgsetter := ""
	if aorg == "a" {
		aorgsetter = "Artist..."
	} else {
		aorgsetter = "Genre..."
	}
	label.SetInnerHTML(aorgsetter)
	adder := d.GetElementByID("moreinput")
	div.AppendChild(input)
	div.AppendChild(label)
	form.AppendChild(div)
	adder.AppendChild(form)
	js.Global.Get("window").Get("componentHandler").Call("upgradeDom")
}

func remove() {
	d := dom.GetWindow().Document()
	d.GetElementsByClassName("inpt")[len(d.GetElementsByClassName("inpt"))-1].SetOuterHTML("")
	if len(d.GetElementsByClassName("inpt")) > 0 {
		d.GetElementByID("remove").(*dom.HTMLButtonElement).Disabled = false
	} else {
		d.GetElementByID("remove").(*dom.HTMLButtonElement).Disabled = true
	}
}

func inptErr() {
	js.Global.Call("alert", "Invalid Input!")
	//js.Global.Get("window").Get("Document").Get("MaterialSnackbar").Call("showSnackbar", "Invalid Input!")
}

//When create button is clicked
func onCreate() {
	d := dom.GetWindow().Document()
	artists := d.GetElementsByClassName("artist")
	overview := d.GetElementByID("overview")
	results := d.GetElementByID("results")
	if len(artists) > 1 {
		all := ""
		for i, _ := range artists {
			val := d.GetElementsByClassName("artist")[i].(*dom.HTMLInputElement).Value
			if i == 0 && all != "" {
				all += val
			} else if i > 0 && all != "" {
				all += "," + val
			}
		}

		artists := aorg + all
		if artists == "a" {
			inptErr()
		} else {
			fmt.Println(artists)
			conn.Write([]byte(artists))
			overview.SetAttribute("style", "display: none;")
			results.SetAttribute("style", "display: block;")
		}
	} else {
		artist := aorg + d.GetElementsByClassName("artist")[0].(*dom.HTMLInputElement).Value
		fmt.Println(artist)
		if artist == "a" {
			inptErr()
		} else {
			fmt.Println(artist)
			conn.Write([]byte(artist))
			overview.SetAttribute("style", "display: none;")
			results.SetAttribute("style", "display: block;")
		}
	}
}
