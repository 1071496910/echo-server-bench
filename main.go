package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	concurrence int
	throughput  int
	server      string
	port        string
)

func init() {
	flag.IntVar(&concurrence, "c", 1, "concurrence num")
	flag.IntVar(&throughput, "t", 1000, "total throughput")
	flag.StringVar(&server, "s", "127.0.0.1", "server addr")
	flag.StringVar(&port, "p", "80", "server port")
}

func main() {
	flag.Parse()

	var n int64
	var wg sync.WaitGroup
	var readBuf = make([]byte, 5)

	begin := time.Now()
	wg.Add(concurrence)
	for i := 0; i < concurrence; i++ {
		go func() {
			defer wg.Done()
			conn, err := net.Dial("tcp", server+":"+port)
			if err != nil {
				panic(err)
			}
			for {
				if atomic.LoadInt64(&n) >= int64(throughput) {
					return
				}
				_, err = conn.Write([]byte("hello"))
				if err != nil {
					panic(err)
				}
				_, err = conn.Read(readBuf)
				if err != nil {
					panic(err)
				}
				atomic.AddInt64(&n, 1)
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(begin)
	fmt.Println("Use time:", elapsed)
}
