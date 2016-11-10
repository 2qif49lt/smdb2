package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
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
	ERR_AUTH
	ERR_SWITCH_MODE
	ERR_RESTART
)

var ERRMSG = map[int]string{
	SUCC:            "成功",
	FAIL:            "失败",
	ERR_FORMAT:      "参数有误",
	ERR_RANGE:       "选择范围有误",
	ERR_INTERNAL:    "内部错误",
	ERR_AUTH:        "授权非法",
	ERR_SWITCH_MODE: "切换模式失败",
	ERR_RESTART:     "重启服务失败",
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

type sendTmpl struct {
	StartTimeMilloSec int
	CurrentMode       string
	ToMode            string
	Msg               string
}

const (
	emailTmplText = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml">
　<head>
　　<meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
　　<title>{{.Title}}</title>
　　<meta name="viewport" content="width=device-width, initial-scale=1.0"/>
 </head>
 <body style="margin: 0; padding: 0;">
　 <table border="0" cellpadding="0" cellspacing="0" width="100%">
　　 <tr> 
　　　 <td> <a href="{{.Url}}" target="_blank" >{{.Url}}</a><br>the one-time address will become invalid after 10 minutes
        <br>@2qif49lt
      </td>
　　 </tr>
　 </table>
  </body>
</html>`
)

type emailTmplStuct struct {
	Title string
	Url   string
}

var emailTmpl, _ = template.New("emailtext").Parse(emailTmplText)

func adminSendUrlHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	email := r.FormValue("email")

	data := sendTmpl{
		StartTimeMilloSec: int(srvStartTime.Unix() * 1000),
	}
	currentmode, err := getDb2Mode()
	tomode := UNKNOWN
	if err != nil {
		currentmode = UNKNOWN
		data.Msg = fmt.Sprintf("check you app.xml: %s", err.Error())
	} else {
		if currentmode == NORMAL {
			tomode = STANDBY
		} else {
			tomode = NORMAL
		}
	}
	data.CurrentMode = currentmode
	data.ToMode = tomode

	if data.Msg == "" && r.Method == "POST" && email != "" && emailTmpl != nil {
		if checkAdminMail(email) == true {
			id := ssnmgr.NewId()
			fmt.Println("new id:", id)
			para := emailTmplStuct{}
			para.Title = "authorised address for switch operation"
			scheme := "http"
			if r.URL.Scheme != "" {
				scheme = r.URL.Scheme
			}
			para.Url = fmt.Sprintf("%s://%s/auth/switch/%s", scheme, r.Host, id)
			sb := bytes.NewBufferString("")
			err := emailTmpl.Execute(sb, para)
			if err == nil {
				content := sb.String()
				err = sendMail2("smdb2", email, fmt.Sprintf("mas操作邮件-切换到%s-%d", tomode, time.Now().Unix()), content)
			}

			if err == nil {
				data.Msg = fmt.Sprintf("OK %s,Pls check your email box!", email)
			} else {
				data.Msg = fmt.Sprintf(`FAIL,%s.check your smtp config!`, err.Error())
			}
		}
	}
	if emailTmpl == nil {
		data.Msg = "FAIL, check you email content template!"
	}

	if tmpl, err := template.ParseFiles("tmpl/send.html"); err == nil {
		tmpl.Execute(w, data)
	} else {
		w.Write([]byte(err.Error()))
	}
}

func authUrlHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	wf := flushWriter{w: w}
	if f, ok := w.(http.Flusher); ok {
		wf.f = f
	}

	var err error = nil

	vars := mux.Vars(r)
	aid := vars["authid"]

	if ssnmgr.IsExist(aid) == false {
		wf.Send("authorization invalid")
		return
	}
	ssnmgr.Del(aid)

	err = switchDb2Mode(wf)
	if err != nil {
		return
	}
	restartDb2Service(wf)

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
	r.HandleFunc("/auth/switch/{authid}", authUrlHandler)
	r.HandleFunc("/admin/send.go", adminSendUrlHandler)
	r.NewRoute().PathPrefix("/admin/").Handler(
		http.StripPrefix("/admin/", http.FileServer(http.Dir("static"))))
	return r
}
