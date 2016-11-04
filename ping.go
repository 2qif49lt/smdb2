package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type response struct {
	addr *net.IPAddr
	rtt  time.Duration
}

func pingRoution(tars []*net.IPAddr, payload int) {
	p := fastping.NewPinger()
	p.Size = payload

	results := make(map[string]*response)
	for _, addr := range tars {
		p.AddIPAddr(addr)
		results[addr.String()] = nil
	}

	onRecv, onIdle := make(chan *response), make(chan bool)

	p.OnRecv = func(addr *net.IPAddr, t time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}
	p.OnIdle = func() {
		onIdle <- true
	}

	p.RunLoop()
	defer p.Stop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	defer signal.Stop(c)

	for {
		select {
		case <-c:
			fmt.Println("get interrupted")
			return
		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:
			now := time.Now()

			for host, r := range results {
				if r == nil {
					pingqueue <- &pingrsp{now, host, -1}
				} else {
					pingqueue <- &pingrsp{now, host, int(r.rtt / time.Millisecond)}
				}
				results[host] = nil
			}
		case <-p.Done():
			if err := p.Err(); err != nil {
				fmt.Println("Ping failed:", err)
			}
			return
		}
	}
}
