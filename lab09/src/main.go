package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"sort"
	"sync"
	"time"
)

func portscan(ip string, port int, timeout time.Duration, wg *sync.WaitGroup, results chan int) {
	defer wg.Done()

	target := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", target, timeout)
	if err != nil {
		return
	}
	conn.Close()
	results <- port
}

func main() {
	ip := flag.String("ip", "", "IP-address to run scanner on")
	startPort := flag.Int("start", 1, "Start port")
	endPort := flag.Int("end", 1024, "End port")
	timeoutSec := flag.Int("timeout", 1, "Timeout")
	workers := flag.Int("workers", 100, "Number of workers to run")
	flag.Parse()

	if *ip == "" {
		flag.Usage()
		return
	}

	if *startPort < 1 || *startPort > 65535 {
		log.Fatal("Bad start port. Must be in range 1 -- 65535")
	}

	if *endPort < 1 || *endPort > 65535 {
		log.Fatal("Bad end port. Must be in range 1 -- 65535")
	}

	if *startPort > *endPort {
		log.Fatal("start port > end port ???")
	}

	timeout := time.Duration(*timeoutSec) * time.Second

	sem := make(chan struct{}, *workers)
	results := make(chan int)
	var wg sync.WaitGroup

	fmt.Printf("Scanning %s\n", *ip)

	for port := *startPort; port <= *endPort; port++ {
		wg.Add(1)
		sem <- struct{}{}
		go func(p int) {
			portscan(*ip, p, timeout, &wg, results)
			<-sem
		}(port)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var openPorts []int
	for port := range results {
		openPorts = append(openPorts, port)
	}
	sort.Ints(openPorts)

	if len(openPorts) > 0 {
		fmt.Println("\nOpened ports:")
		for _, port := range openPorts {
			fmt.Printf("%d\n", port)
		}
	} else {
		fmt.Println("No available ports")
	}
}
