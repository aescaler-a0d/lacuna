/*
 * @File: arp.go
 * @Date: 2019-05-30 17:47:33
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-03 18:04:09
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"errors"
	"fmt"
)

// Sends hosts in ArpHosts to output chan
//func ArpHosts(done <-chan struct{}, ips <-chan net.IP, c chan<- net.IP, arp []net.IP) {
func ArpHosts(w workGenerator) error {
	// for every IP that comes across the channel
	for ip := range w.i { // HLpaths
		// for every IP in the filter
		if w.debug {
			fmt.Printf("ArpHosts worker: Processing IP: %v\n", ip)
		}
		for _, v := range w.filter {
			// check if the IP in the channel is in the filter
			if ip.Equal(v) {
				select {
				case w.o <- ip:
					if w.debug {
						fmt.Printf("ArpHosts worker: IP in filter: %v\n", ip)
					}
				case <-w.done:
					if w.debug {
						fmt.Printf("ArpHosts worker: Done: called inside if\n")
					}
					return nil
				}
			}
			select {
			case <-w.done:
				if w.debug {
					fmt.Printf("ArpHosts worker: Done: called outside if\n")
				}
				return nil
			}
		}
	}
	return errors.New("ArpHosts worker: Error: out of range before calling done")
}
