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

type pingrsp struct {
	T   time.Time `json:"t"`
	Tar string    `json:"tar"`
	Ms  int       `json:"ms"`
}

var pingqueue = make(chan *pingrsp, 100)

const (
	timekeyformat = "20060102150405"
)
