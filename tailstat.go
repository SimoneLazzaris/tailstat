package main

import (
	"bufio"
	"time"
	"fmt"
	"strings"
	"sync"
	"net"
)


var mtx=&sync.Mutex{}

func scanMetrics(s string) {
	mtx.Lock()
	for i, m:= range Cfg.Metrics {
		if strings.Contains(s, m.Match) {
			m.Counter++
			Cfg.Metrics[i]=m
		}
	}
	mtx.Unlock()
}

func printMetrics() {
	t:=time.Now()
	for {
		t=t.Add(time.Duration(Cfg.Deltat)*time.Second)
		tnow:=time.Now()
		time.Sleep(t.Sub(tnow))
		mtx.Lock()
		for i, m:= range Cfg.Metrics {
			//fmt.Printf("%s: %d ", m.Name, m.Counter)
			sendMetric(tnow, m.Name, m.Counter)
			m.Counter=0
			Cfg.Metrics[i]=m
		}
		mtx.Unlock()
		//fmt.Println("")
	}
}

func sendMetric(t time.Time, name string, Counter int) {
	conn, err := net.Dial("tcp", Cfg.Graphite)
	if err!=nil {
		return
	}
	s:=fmt.Sprintf("%s.%s %d %d\n",Cfg.Prefix,name,Counter,t.Unix())
	bb:=bufio.NewWriter(conn)
	bb.WriteString(s)
	bb.Flush()
// 	fmt.Printf("%s.%s %d %d\n",Cfg.Prefix,name,Counter,t.Unix())
}

func main() {
// 	fmt.Println("starting")
	initCfg()
//	daemon(0,0)
// 	fmt.Println("started")
	go printMetrics()
	ff:=follower{}
	ff.init(Cfg.Logfile)
	for {
		scanMetrics(ff.tail())
	}
}
