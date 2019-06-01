/*
 * @File: dns.go
 * @Date: 2019-05-30 17:48:02
 * @OA:   antonioe
 * @CA:   antonioe
 * @Time: 2019-06-01 00:27:02
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"
	"sync"
	"time"
)

func DnsHosts(s *net.IPNet) ([]net.IP, error) {
	return DnsHostsAsync(s)
}

func DnsHostsAsync(s *net.IPNet) ([]net.IP, error) {
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
				// var err error
				// pinger, err := ping.NewPinger(ip.String())
				// if err != nil {
				// 	fmt.Println("Error: Could not create pinger")
				// }
				// pinger.Count = attempts
				// pinger.Interval = interval
				// pinger.Timeout = timeout
				// pinger.SetPrivileged(true)
				// pinger.Run()
				// if pinger.PacketsSent > 0 && pinger.PacketsRecv > 0 {
				// 	res <- ip
				// }
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
