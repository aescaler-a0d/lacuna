/*
 * @File: arpwatch.go
 * @Date: 2019-05-30 17:47:47
 * @OA:   antonioe
 * @CA:   Antonio Escalera
 * @Time: 2019-06-01 19:48:19
 * @Mail: antonioe@wolfram.com
 * @Copy: Copyright © 2019 Antonio Escalera <aj@angelofdeauth.host>
 */

package find

import (
	"net"
	"time"

	"github.com/angelofdeauth/gopher/read"
)

func ArpWatch(s *net.IPNet, a string, debug bool) ([]net.IP, error) {
	n := []net.IP{}
	awd, err := read.ReadAWDataIntoStruct(a)
	if err != nil {
		return n, err
	}
	t := time.Now().Unix()
	for _, v := range awd {
		if t-v.Time < 15552000 {
			n = append(n, v.Ipaddr)
		}
	}
	return n, nil
}
