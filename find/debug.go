/*
 * @File: debug.go
 * @Date: 2019-05-30 17:47:33
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-06 16:56:40
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"fmt"
	"net"
)

// Sends hosts in ArpHosts to output chan
func Debug(w workGenerator, ip net.IP) net.IP {
	if w.debug {
		fmt.Printf("Debug worker: Processing IP: %v\n", ip)
	}
	return ip
}
