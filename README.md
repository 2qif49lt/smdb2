# smdb2
monitor db2 service,use for huawei mas device.


####usage

1. go get github.com/2qif49lt/smdb2
2. config db2 include, lib path in build.sh and run.sh
3. build
4. config the command line parameters, config.toml especially
5. run 

####notice
1. *need ROOT privilege*
2. you may consider to config the port in iptables for better security 

####web interface information
as follows
```go
    r.HandleFunc("/db2/ws.go", wsDB2Handler)
    r.HandleFunc("/ping/ws.go", wsPingHandler)
    r.HandleFunc("/chart/ping.go", pingChartHandler)
    r.HandleFunc("/chart/db2.go", dbChartHandler)
    r.HandleFunc("/ping/last", pingLastHandler)
    r.HandleFunc("/ping/{tar}", pingHandler)
    r.HandleFunc("/db2/status", db2StatusHandler)
    r.HandleFunc("/auth/switch/{authid}", authUrlHandler)
    r.HandleFunc("/admin/send.go", adminSendUrlHandler)
```