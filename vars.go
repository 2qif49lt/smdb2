package main

import (
	"time"
)

type dbcount struct {
	T       time.Time `json:"t"`
	Sending int       `json:"sending"`
	Ok      int       `json:"ok"`
	Fail    int       `json:"fail"`
}

const (
	timekeyformat = "200601021504"
	boltdbname    = "sm.db"
	version       = "1.0.0"
)

var (
	cntqueue        = make(chan *dbcount, 100)
	pingReduceQueue = make(chan *pingReduceRsp, 100)
	ssnmgr          = newssnMgr()
	srvStartTime    = time.Now()
)

type pingReduceRsp struct {
	T time.Time `json:"t"`

	Tar string `json:"tar"`

	Ave  int `json:"ave"`
	Min  int `json:"min"`
	Max  int `json:"max"`
	Send int `json:"send"`
	Rev  int `json:"rev"`
}
