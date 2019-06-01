/*
 * @File: arpwatch.go
 * @Date: 2019-05-30 17:47:47
 * @OA:   antonioe
 * @CA:   antonioe
 * @Time: 2019-05-31 03:13:37
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"
	"time"

	"github.com/angelofdeauth/gopher/read"
)

func ArpWatch(s *net.IPNet, a string) (n []net.IP) {
	awd := read.ReadAWDataIntoStruct(a)
	t := time.Now().Unix()
	for _, v := range awd {
		if t-v.Time < 15552000 {
			n = append(n, v.Ipaddr)
		}
	}
	return n
}
