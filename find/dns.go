/*
 * @File: dns.go
 * @Date: 2019-05-30 17:48:02
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-01 20:26:12
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"fmt"
	"net"
	"sync"
)

func DnsHosts(s *net.IPNet, debug bool) ([]net.IP, error) {
	return DnsHostsAsync(s, debug)
}

func DnsHostsAsync(s *net.IPNet, debug bool) ([]net.IP, error) {
	// parse CIDR arguments
	sl := []net.IP{}
	generator, err := HostsInSubnet(s)
	if err != nil {
		return sl, err
	}

	// prepare worker
	wg := &sync.WaitGroup{}
	wg.Add(poolSize)
	ips := make(chan net.IP, poolSize)
	res := make(chan net.IP, poolSize)

	for i := 0; i < poolSize; i++ {
		go func(thread int) {
			for ip := range ips {
				if debug {
					fmt.Printf("IP: %v THREAD: %v\n", ip, thread)
				}
				hname, _ := net.LookupAddr(ip.String())
				if len(hname) > 0 {
					for _, v := range hname {
						if v != "" {
							res <- ip
						}
					}
				}
			}
			wg.Done()
		}(i)
	}

	// printer
	// pr := &sync.WaitGroup{}
	// pr.Add(1)
	// go func() {
	// 	bar := pb.New64(int64(total))
	// 	bar.ShowBar = true
	// 	bar.ShowTimeLeft = false
	// 	bar.ShowCounters = true
	// 	bar.Start()
	// 	const clear = "\x1b[2K\r" // ansi delete line + CR
	// 	for ip := range ips {
	// 		bar.Increment()
	// 		if ip != nil {
	// 			log.Printf("%s DNS: %s", clear, ip)
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
			return nil
		})
	}

	// wait for worker and printer to finish
	close(ips)
	wg.Wait()
	close(res)
	sl = ChanToSlice(res).([]net.IP)

	return sl, nil
}
