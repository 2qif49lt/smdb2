package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

type httprsp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

const (
	SUCC int = iota
	FAIL
	ERR_FORMAT
	ERR_RANGE
	ERR_INTERNAL
)

var ERRMSG = map[int]string{
	SUCC:         "成功",
	FAIL:         "失败",
	ERR_FORMAT:   "参数有误",
	ERR_RANGE:    "选择范围有误",
	ERR_INTERNAL: "内部错误",
}

func smHandler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "smhandler")
}

func pingHandler(w http.ResponseWriter, req *http.Request) {
	msg := &httprsp{}
	code := FAIL
	errdesc := ""
	var data interface{} = nil
	var err error = nil

	vars := mux.Vars(req)
	tar := vars["tar"]

	req.ParseForm()

	from := req.FormValue("from")
	to := req.FormValue("to")
	slimit := req.FormValue("limit")
	limit := 0
	if slimit != "" {
		limit, err = strconv.Atoi(slimit)
		if err != nil {
			errdesc = err.Error()
			code = ERR_FORMAT
			goto END
		}
		if limit < 0 {
			errdesc = "limit只能是大于等于0"
			code = ERR_FORMAT
			goto END
		}
	}
	if limit == 0 {
		limit = 365 * 24 * 60
	}

	if from == "" {
		from = time.Now().Add(time.Minute * -60).Format(timekeyformat)
	}

	if to == "" {
		to = time.Now().Format(timekeyformat)
	}

	if from > to {
		code = ERR_RANGE
		goto END
	}

	err, data = dbReadPing(tar, from, to, limit)
	if err != nil {
		errdesc = err.Error()
		code = ERR_INTERNAL
		goto END
	}
	code = SUCC
END:
	msg.Code = code
	msg.Msg = ERRMSG[msg.Code]
	if errdesc != "" {
		msg.Msg += errdesc
	}
	if msg.Code == SUCC {
		msg.Data = data
	}
	ret, _ := json.Marshal(msg)

	w.Write(ret)
}

func pingLastHandler(w http.ResponseWriter, req *http.Request) {
	msg := &httprsp{}
	code := FAIL
	errdesc := ""
	var data interface{} = nil
	var err error = nil

	req.ParseForm()

	slimit := req.FormValue("limit")
	limit := 0

	if slimit != "" {
		limit, err = strconv.Atoi(slimit)
		if err != nil {
			errdesc = err.Error()
			code = ERR_FORMAT
			goto END
		}
		if limit <= 0 {
			errdesc = "limit只能是大于等于0"
			code = ERR_FORMAT
			goto END
		}
	}
	if limit == 0 {
		limit = 10
	}

	err, data = dbReadPingLast(limit)
	if err != nil {
		errdesc = err.Error()
		code = ERR_INTERNAL
		goto END
	}
	code = SUCC
END:
	msg.Code = code
	msg.Msg = ERRMSG[msg.Code]
	if errdesc != "" {
		msg.Msg += errdesc
	}
	if msg.Code == SUCC {
		msg.Data = data
	}
	ret, _ := json.Marshal(msg)

	w.Write(ret)
}

func db2StatusHandler(w http.ResponseWriter, req *http.Request) {
	msg := &httprsp{}
	code := FAIL
	errdesc := ""
	var data interface{} = nil
	var err error = nil

	req.ParseForm()

	from := req.FormValue("from")
	to := req.FormValue("to")
	slimit := req.FormValue("limit")
	limit := 0
	if slimit != "" {
		limit, err = strconv.Atoi(slimit)
		if err != nil {
			errdesc = err.Error()
			code = ERR_FORMAT
			goto END
		}
		if limit < 0 {
			errdesc = "limit只能是大于等于0"
			code = ERR_FORMAT
			goto END
		}
	}
	if limit == 0 {
		limit = 365 * 24 * 60
	}

	if from == "" {
		from = time.Now().Add(time.Minute * -60).Format(timekeyformat)
	}

	if to == "" {
		to = time.Now().Format(timekeyformat)
	}

	if from > to {
		code = ERR_RANGE
		goto END
	}

	err, data = dbReadDb2(from, to, limit)
	if err != nil {
		errdesc = err.Error()
		code = ERR_INTERNAL
		goto END
	}
	code = SUCC
END:
	msg.Code = code
	msg.Msg = ERRMSG[msg.Code]
	if errdesc != "" {
		msg.Msg += errdesc
	}
	if msg.Code == SUCC {
		msg.Data = data
	}
	ret, _ := json.Marshal(msg)

	w.Write(ret)
}

func httpSrv(port int) error {
	http.Handle("/", regRouter())

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func regRouter() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "sm monitor %s \n%s", version,
			time.Now().Format(time.RFC3339))
	})

	r.HandleFunc("/sm", smHandler)
	r.HandleFunc("/ping/last", pingLastHandler)
	r.HandleFunc("/ping/{tar}", pingHandler)
	r.HandleFunc("/db2/status", db2StatusHandler)

	return r
}
