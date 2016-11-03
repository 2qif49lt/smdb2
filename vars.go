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

var cntqueue = make(chan dbcount, 100)

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

var pingqueue = make(chan *response, 100)
