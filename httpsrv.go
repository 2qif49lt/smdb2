package main

import (
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
	"time"
)

type httprsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	SUCC int = itoa
	FAIL
	ERR_FORMAT
)

var ERRMSG = map[int]string{
	SUCC:       "成功",
	FAIL:       "失败",
	ERR_FORMAT: "参数格式有误",
}

func smHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "smhandler")
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
	msg := &httprsp{}
	code := FAIL
	var data interface{} = nil

	vars := mux.Vars(r)
	tar := vars["tar"]

	req.ParseForm()

	from := req.FormValue("from")
	to := req.FormValue("to")

	if from == "" {
		from = time.Now().Add(time.Minute * -30).Format(timekeyformat)
	} else {
		iaz := len(timekeyformat) - len(from)
		if iaz > 0 {
			from += strings.Repeat("0", iaz)
		}

		if _, err := time.Parse(timekeyformat, from); err != nil {
			code = ERR_FORMAT
			goto END
		}
	}

	if to == "" {
		to = time.Now().Format(timekeyformat)
	} else {
		iaz := len(timekeyformat) - len(from)
		if iaz > 0 {
			from += strings.Repeat("0", iaz)
		}

		if _, err := time.Parse(timekeyformat, from); err != nil {
			code = ERR_FORMAT
			goto END
		}
	}

END:
	msg.Code = code
	msg.Msg = ERRMSG[msg.Code]
	if msg.Code == SUCC {
		msg.Data = data
	}
	jsonmsg, _ := json.Marshal(msg)

	w.Write(ret)
}
func db2Handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "db2Handler")
}

func httpSrv(port int) error {
	http.Handle("/", regRouter())

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func dbReadPing(db *bolt.DB, tar string, from, to string) error {
	err := db.View(func(tx *bolt.Tx) error {
		bt := tx.Bucket([]byte(fmt.Sprintf(`ping-%s`, tar)))
		if bt == nil {
			return fmt.Errorf("ip donot exist")
		}
		cr := bt.Cursor()
		if cr == nil {
			return fmt.Errorf("cursor return nil")
		}
		fkey, fvalue := cr.Seek([]byte(from))
		if fkey == nil {
			return fmt.Errorf("%s key is not exist,and no keys follow", from)
		}

		for k, v := c.Seek([]byte(from)); k != nil && bytes.Compare(k, max) <= 0; k, v = c.Next() {
			fmt.Printf("%s: %s\n", k, v)
		}
		return nil
	})
	return nil
}

func regRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "sm monitor %s", time.Now().Format(time.RFC3339))
	})

	r.HandleFunc("/sm", smHandler)
	r.HandleFunc("/ping/{tar}", pingHandler)
	r.HandleFunc("/db2", db2Handler)

	return r
}
