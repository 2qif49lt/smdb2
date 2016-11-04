package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	host                   = flag.String("host", "", "host to connect")
	port                   = flag.Int("port", 50110, "port to connect")
	usrname                = flag.String("user", "", "user name")
	password               = flag.String("pwd", "", "password")
	dbname                 = flag.String("db", "", "database to connect")
	httpport               = flag.Int("srvport", 12345, "srv port for listening")
	tarip                  = flag.String("tarip", "", "target ip to ping,split by ,")
	payload                = flag.Int("pingload", 1024, "size in bytes of the payload to ping, at least 8")
	tars     []*net.IPAddr = nil
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: %s [options]

%s need set parameters to connect to DB2 

example:

 %s -host=127.0.0.1 -port=xxxx -user=xxxx -pwd=xxxx -db=xxxxx
`, os.Args[0], os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func checkTars(tars string) []*net.IPAddr {
	naddr := make([]*net.IPAddr, 0)
	addrs := strings.Split(tars, ",")

	if len(addrs) == 0 {
		return nil
	}
	for _, addr := range addrs {
		ra, err := net.ResolveIPAddr("ip4:icmp", addr)
		if err != nil {
			return nil
		}
		naddr = append(naddr, ra)
	}
	return naddr
}
func checkArge() {
	flag.Usage = usage
	flag.Parse()

	tars = checkTars(*tarip)

	if *dbname == "" || *host == "" || *port == 0 ||
		*usrname == "" || *password == "" || *dbname == "" ||
		*httpport == 0 || tars == nil || *payload < 8 {
		flag.Usage()
	}
}

func main() {
	checkArge()

	conn := fmt.Sprintf(`DATABASE=%s; HOSTNAME=%s; PORT=%d; PROTOCOL=TCPIP; UID=%s; PWD=%s;`,
		*dbname, *host, *port, *usrname, *password)

	go db2Roution(conn)
	go boltWriteRoution()
	go pingRoution(tars, payload)

	err := httpSrv(*httpport)
	if err != nil {
		fmt.Println("httpSrv: ", err)
	}
}
