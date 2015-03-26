package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"github.com/tom-f/brouker/broker"
)

var (
	addr = flag.String("addr", ":8050", "http service address")
)

func homeHandler(c http.ResponseWriter, req *http.Request) {
	tpl := template.Must(template.ParseFiles("./home.html"))
	err := tpl.Execute(c, req.Host)

	if err != nil {
		fmt.Println(err)
	}
	// homeTempl = template.Must(template.ParseFiles(filepath.Join(*assets, "./home.html")))

}

func main() {
	flag.Parse()
	c := broker.NewCtrl()
	go c.Run()
	http.HandleFunc("/", homeHandler)
	http.Handle("/ws", broker.ConnHandler{C: c})
	http.Handle("/msg", broker.CtrlWriter{C: c})
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
