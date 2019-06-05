/*
 * @File: arpwatch.go
 * @Date: 2019-05-30 17:47:47
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-05 13:02:53
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"fmt"
	"net"
)

// Sends hosts in ArpWatch to output chan.
func ArpWatch(w workGenerator, ip net.IP) net.IP {
	if w.debug {
		fmt.Printf("ArpWatch worker: Processing IP: %v\n", ip)
	}
	for _, v := range w.filter {
		// check if the ip in the channel is in the filter
		if ip.Equal(v) {
			if w.debug {
				fmt.Printf("ArpWatch worker: IP in filter: %v\n", ip)
			}
			return ip
		}
	}
	return nil
}
