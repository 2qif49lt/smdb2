package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
)

var (
	host                   = flag.String("host", "127.0.0.1", "db host to connect")
	port                   = flag.Int("port", 50110, "db port to connect")
	usrname                = flag.String("user", "", "db user name")
	password               = flag.String("pwd", "", "db user'spassword")
	dbname                 = flag.String("db", "", "database to connect")
	httpport               = flag.Int("srvport", 12345, "this srv's port for listening")
	tarips                 = flag.String("tarips", "", "target ip to ping,split by ,")
	payload                = flag.Int("pingload", 1024, "size in bytes of the payload to ping, at least 8")
	cfgpath                = flag.String("cfgpath", "config.toml", "config file path")
	tars     []*net.IPAddr = nil
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: %s [options]

need set parameters to connect to DB2 

example:

 %s -host=127.0.0.1 -port=1234 -user=who -pwd=what -db=where
`, os.Args[0], os.Args[0])
	flag.PrintDefaults()
	os.Exit(1)
}

func checkTars(tarips string) ([]*net.IPAddr, error) {
	if tarips == "" {
		return nil, fmt.Errorf("tarips is empty")
	}

	naddr := make([]*net.IPAddr, 0)
	addrs := strings.Split(tarips, ",")

	if len(addrs) == 0 {
		return nil, fmt.Errorf("tarips is empty")
	}

	for _, addr := range addrs {
		ra, err := net.ResolveIPAddr("ip4:icmp", addr)
		if err != nil {
			return nil, err
		}
		naddr = append(naddr, ra)
	}
	return naddr, nil
}

func checkConfig() error {
	cfg = &ConfigFile{Filename: *cfgpath}
	return cfg.Load()
}

func checkArge() {
	flag.Usage = usage
	flag.Parse()
	var err error = nil
	tars, err = checkTars(*tarips)

	if *dbname == "" || *host == "" || *port == 0 ||
		*usrname == "" || *password == "" || *dbname == "" ||
		*httpport == 0 || tars == nil || err != nil || *payload < 8 {
		flag.Usage()
	}
}

func main() {
	checkArge()

	if err := checkConfig(); err != nil {
		fmt.Println(err)
		return
	}
	conn := fmt.Sprintf(`DATABASE=%s; HOSTNAME=%s; PORT=%d; PROTOCOL=TCPIP; UID=%s; PWD=%s;`,
		*dbname, *host, *port, *usrname, *password)

	go boltWriteRoution()
	go db2Roution(conn)
	go pingRoution(tars, *payload)
	go sessionAliveRountion()

	err := httpSrv(*httpport)
	if err != nil {
		fmt.Println("httpSrv: ", err)
	}
}
