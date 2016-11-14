package main

import (
	"encoding/json"
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

	err, rsps := dbReadPingLastHour(60)
	if err == nil {
		for idx := len(rsps) - 1; idx >= 0; idx-- {
			err = c.WriteMessage(websocket.TextMessage, []byte(rsps[idx]))
			if err != nil {
				fmt.Println("dbReadPingLastHour WriteMessage", err)
				break
			}
		}
	}
	sub := puber.SubscribeTopic(func(v interface{}) bool {
		_, ok := v.(*pingReduceRsp)
		return ok
	})
	defer puber.Evict(sub)

	for {
		select {
		case data := <-sub:
			pingrsp := data.(*pingReduceRsp)

			val, err := json.Marshal(pingrsp)
			if err != nil {
				fmt.Println("marshl", err)
				return
			}

			err = c.WriteMessage(websocket.TextMessage, val)
			if err != nil {
				fmt.Println("writemessage", err)
				return
			}
		default:
			_, _, err := c.ReadMessage()
			if err != nil {
				fmt.Println("read:", err)
				return
			}
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
