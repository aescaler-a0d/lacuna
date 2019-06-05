/*
 * @File: arp.go
 * @Date: 2019-05-30 17:47:33
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-05 13:02:45
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"fmt"
	"net"
)

// Sends hosts in ArpHosts to output chan
func ArpHosts(w workGenerator, ip net.IP) net.IP {
	if w.debug {
		fmt.Printf("ArpHosts worker: Processing IP: %v\n", ip)
	}
	for _, v := range w.filter {
		// check if the IP in the channel is in the filter
		if ip.Equal(v) {
			if w.debug {
				fmt.Printf("ArpHosts worker: IP in filter: %v\n", ip)
			}
			return ip
		}
	}
	return nil
}
