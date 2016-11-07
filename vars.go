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

var cntqueue = make(chan *dbcount, 100)

const (
	timekeyformat = "200601021504"
	boltdbname    = "sm.db"
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

var pingReduceQueue = make(chan *pingReduceRsp, 100)

const (
	version = "1.0.0"
)
