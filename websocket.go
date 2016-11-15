package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type wsnum struct {
	Num int
	Ws  string
}

func pingChartHandler(w http.ResponseWriter, r *http.Request) {
	data := wsnum{
		30,
		fmt.Sprintf("ws://%s/ping/ws.go", r.Host),
	}
	r.ParseForm()
	num := r.FormValue("num")
	if num != "" {
		num := r.FormValue("num")
		if i, e := strconv.Atoi(num); e == nil {
			data.Num = i
		}
	}

	if tmpl, err := template.ParseFiles("tmpl/ping.html"); err == nil {
		tmpl.Execute(w, data)
	} else {
		w.Write([]byte(err.Error()))
	}
}

func dbChartHandler(w http.ResponseWriter, r *http.Request) {
	data := wsnum{
		30,
		fmt.Sprintf("ws://%s/db2/ws.go", r.Host),
	}
	r.ParseForm()
	num := r.FormValue("num")
	if num != "" {
		if i, e := strconv.Atoi(num); e == nil {
			data.Num = i
		}
	}

	if tmpl, err := template.ParseFiles("tmpl/db.html"); err == nil {
		tmpl.Execute(w, data)
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

	c.WriteMessage(websocket.TextMessage, []byte("hi"))

	_, numbyte, err := c.ReadMessage()
	if err != nil {
		fmt.Print("ReadMessage:", err)
		return
	}
	num, err := strconv.Atoi(string(numbyte))
	if err != nil || num <= 0 {
		return
	}

	err, rsps := dbReadPingLastDay(num)
	if err == nil {
		for idx := len(rsps) - 1; idx >= 0; idx-- {
			err = c.WriteMessage(websocket.TextMessage, []byte(rsps[idx]))
			if err != nil {
				fmt.Println("dbReadPingLastDay WriteMessage", err)
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

	c.WriteMessage(websocket.TextMessage, []byte("hi"))

	_, numbyte, err := c.ReadMessage()
	if err != nil {
		fmt.Print("ReadMessage:", err)
		return
	}
	num, err := strconv.Atoi(string(numbyte))
	if err != nil || num <= 0 {
		return
	}

	err, rsps := dbReadDb2LastDay(num)
	if err == nil {
		for idx := len(rsps) - 1; idx >= 0; idx-- {
			err = c.WriteMessage(websocket.TextMessage, []byte(rsps[idx]))
			if err != nil {
				fmt.Println("dbReadDb2LastDay WriteMessage", err)
				break
			}
		}
	}
	sub := puber.SubscribeTopic(func(v interface{}) bool {
		_, ok := v.(*dbcount)
		return ok
	})
	defer puber.Evict(sub)

	for {
		select {
		case data := <-sub:
			cnt := data.(*dbcount)

			val, err := json.Marshal(cnt)
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
