/*
 * @File: dns.go
 * @Date: 2019-06-07 12:56:20
 * @OA:   Antonio Escalera
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 13:28:32
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */
package read

import (
	"net"
)

func DnsHosts(s *net.IPNet, debug bool) ([]net.IP, error) {
	return DnsHostsAsync(s, debug)
}
func DnsHostsAsync(s *net.IPNet, debug bool) ([]net.IP, error) {
	// parse CIDR arguments
	sl := []net.IP{}
	//generator, err := ReadSubnetIntoChan(s, debug)
	//hname, _ := net.LookupAddr(ip.String())
	//if len(hname) > 0 {
	//	for _, v := range hname {
	//		if v != "" {
	//			res <- ip
	//		}
	//	}
	//}
	return sl, nil
}
