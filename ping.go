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

	for _, addr := range tars {
		p.AddIPAddr(addr)
	}

	results := make(map[string]*response)

	onRecv, onIdle := make(chan *response), make(chan bool)

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		onRecv <- &response{addr: addr, rtt: t}
	}
	p.OnIdle = func() {
		onIdle <- true
	}

	p.RunLoop()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	for {
		select {
		case <-c:
			fmt.Println("get interrupted")
			break
		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:
			now := time.Now()

			for host, r := range results {
				if r == nil {
					fmt.Printf("%s : unreachable %v\n", host, now)
					pingqueue <- &pingrsp{now, host, -1}
				} else {
					fmt.Printf("%s : %v %v\n", host, r.rtt/time.Millisecond, now)
					pingqueue <- &pingrsp{now, host, int(r.rtt / time.Millisecond)}
				}
				results[host] = nil
			}
		case <-p.Done():
			if err = p.Err(); err != nil {
				fmt.Println("Ping failed:", err)
			}
			break
		}
	}
	signal.Stop(c)
	p.Stop()
}
