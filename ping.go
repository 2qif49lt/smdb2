package main

import (
	"fmt"
	"github.com/tatsushid/go-fastping"
	"net"
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
	rsprst := newRspRst()

	for _, addr := range tars {
		p.AddIPAddr(addr)
		results[addr.String()] = nil
		rsprst.AddAddr(addr.String())
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

	tbeg := time.Now()

	for {
		select {
		case <-time.After(time.Minute):

		case res := <-onRecv:
			if _, ok := results[res.addr.String()]; ok {
				results[res.addr.String()] = res
			}
		case <-onIdle:

			for host, r := range results {
				if r == nil {
					rsprst.Append(host, &response{addr: nil, rtt: maxTimeOutTTL * time.Millisecond})
				} else {
					rsprst.Append(host, r)
				}
				results[host] = nil

				if time.Since(tbeg) > time.Minute {
					rsprst.ReduceFor(time.Now(), func(r *pingReduceRsp) {
						pingReduceQueue <- r
					})
					rsprst.Clear()
					tbeg = time.Now()
				}
			}
		case <-p.Done():
			if err := p.Err(); err != nil {
				fmt.Println("Ping failed:", err)
			}
			return
		}
	}
}

const (
	maxTimeOutTTL = -1
)

type rspRst map[string][]*response

func newRspRst() rspRst {
	m := make(map[string][]*response)
	return rspRst(m)
}
func (rst rspRst) AddAddr(addr string) {
	rst[addr] = []*response{}
}

func (rst rspRst) Clear() {
	for k, _ := range rst {
		rst[k] = []*response{}
	}
}

func (rst rspRst) Append(addr string, rsp *response) {
	if _, exist := rst[addr]; exist {
		rst[addr] = append(rst[addr], rsp)
	}
}

func (rst rspRst) ReduceFor(now time.Time, r func(*pingReduceRsp)) {
	for addr, rsps := range rst {
		reduce := &pingReduceRsp{}

		reduce.T = now
		reduce.Tar = addr
		reduce.Send = len(rsps)
		reduce.Max = maxTimeOutTTL
		reduce.Min = maxTimeOutTTL

		for _, rsp := range rsps {
			ms := int(rsp.rtt / time.Millisecond)

			if ms != maxTimeOutTTL {
				reduce.Rev++
				reduce.Ave += ms
			}
			if reduce.Max == maxTimeOutTTL {
				reduce.Max = ms
			}
			if reduce.Min == maxTimeOutTTL {
				reduce.Min = ms
			}

			if ms > reduce.Max {
				reduce.Max = ms
			}

			if ms < reduce.Min {
				reduce.Min = ms
			}
		}

		if reduce.Rev != 0 {
			reduce.Ave /= reduce.Rev
		} else {
			reduce.Ave = maxTimeOutTTL
		}

		r(reduce)
	}
}
