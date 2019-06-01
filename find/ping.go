/*
 * @File: ping.go
 * @Date: 2019-05-30 17:47:27
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-01 13:25:19
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	ping "github.com/sparrc/go-ping"
)

func PingHosts(s *net.IPNet, i string) ([]net.IP, error) {
	inf, err := net.InterfaceByName(i)
	if err != nil {
		return []net.IP{}, err
	}
	a, err := inf.Addrs()
	if err != nil {
		return []net.IP{}, err
	}
	for _, b := range a {
		ip := b.String()
		nip, _, _ := net.ParseCIDR(ip)
		if (strings.Count(ip, ":") < 2) && (s.Contains(nip)) {
			// found first IP on interface i that is in subnet s
			return PingHostsAsync(s, ip, inf.Name)
		}
	}
	// no IP found, return error
	err = errors.New("No assigned address in subnet.")
	return []net.IP{}, err
}

func PingHostsAsync(s *net.IPNet, b string, inf string) ([]net.IP, error) {
	// parse CIDR arguments
	generator, err := HostsInSubnet(s)
	if err != nil {
		return []net.IP{}, err
	}

	// total := len(generator)

	// prepare worker
	wg := &sync.WaitGroup{}
	wg.Add(poolSize)
	ips := make(chan net.IP, poolSize)
	res := make(chan net.IP, poolSize)

	for i := 0; i < poolSize; i++ {
		go func() {
			for ip := range ips {
				var err error
				pinger, err := ping.NewPinger(ip.String())
				if err != nil {
					fmt.Println("Error: Could not create pinger")
				}
				pinger.Count = attempts
				pinger.Interval = interval
				pinger.Timeout = timeout
				pinger.SetPrivileged(true)
				pinger.Run()
				if pinger.PacketsSent > 0 && pinger.PacketsRecv > 0 {
					res <- ip
				}
			}
			wg.Done()
		}()
	}

	// printer
	// pr := &sync.WaitGroup{}
	// pr.Add(1)
	// go func() {
	// 	bar := pb.New64(int64(total))
	// 	bar.ShowBar = true
	// 	bar.ShowTimeLeft = true
	// 	bar.ShowCounters = true
	// 	bar.Start()
	// 	const clear = "\x1b[2K\r" // ansi delete line + CR
	// 	for r := range res {
	// 		bar.Increment()
	// 		if r != nil {
	// 			log.Printf("%s%s", clear, r)
	// 			bar.Update()
	// 		}
	// 	}
	// 	bar.Finish()
	// 	pr.Done()
	// }()

	// yield all IP addresses
	for _, g := range generator {
		each(g, func(ip net.IP) error {
			ips <- ip
			time.Sleep(interval)
			return nil
		})
	}

	// wait for worker and printer to finish
	close(ips)
	wg.Wait()
	close(res)
	sl := ChanToSlice(res).([]net.IP)

	return sl, nil
}
