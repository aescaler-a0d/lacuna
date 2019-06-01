/*
 * @File: arp.go
 * @Date: 2019-05-30 17:47:33
 * @OA:   antonioe
 * @CA:   antonioe
 * @Time: 2019-05-31 04:04:34
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"

	"github.com/angelofdeauth/gopher/read"
)

func ArpHosts(s *net.IPNet, a string) (n []net.IP) {
	arp := read.ReadArpDataIntoStruct(a)
	for _, v := range arp {
		n = append(n, v.Ipaddr)
	}
	return n
}
