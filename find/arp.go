/*
 * @File: arp.go
 * @Date: 2019-05-30 17:47:33
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-01 19:48:12
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"

	"github.com/angelofdeauth/gopher/read"
)

func ArpHosts(s *net.IPNet, a string, debug bool) ([]net.IP, error) {
	n := []net.IP{}
	arp, err := read.ReadArpDataIntoStruct(a)
	if err != nil {
		return n, err
	}
	for _, v := range arp {
		n = append(n, v.Ipaddr)
	}
	return n, nil
}
