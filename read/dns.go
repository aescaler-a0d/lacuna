/*
 * @File: dns.go
 * @Date: 2019-06-07 12:56:20
 * @OA:   Antonio Escalera
 * @CA:   Antonio Escalera
 * @Time: 2019-06-07 14:40:48
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */
package read

import (
	"net"
	"strings"
)

func DnsHostnames(debug bool, ip net.IP) ([]string, error) {
	hns, err := net.LookupAddr(ip.String())
	if err != nil && !strings.Contains(err.Error(), "Name or service not known") {
		return nil, err
	}
	r := make([]string, 0)
	if len(hns) > 0 {
		for _, v := range hns {
			if v != "" {
				r = append(r, v)
			}
		}
	}
	return r, nil
}
