package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/websocket"
)

func pingChartHandler(w http.ResponseWriter, r *http.Request) {
	if tmpl, err := template.ParseFiles("tmpl/ping.html"); err == nil {
		tmpl.Execute(w, "ws://"+r.Host+"/ping/ws.go")
	} else {
		w.Write([]byte(err.Error()))
	}
}

var upgrader = websocket.Upgrader{}

func wsPingHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, _, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}

	}
}

func wsDB2Handler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			fmt.Println("write:", err)
			break
		}
	}
}
